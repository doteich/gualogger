use crate::crd::{GuaLogger, GuaLoggerKubeImage};
use futures::stream::StreamExt;
use kube::{
    Client, Error, Resource,
    api::Api,
    runtime::{Controller, controller::Action, watcher::Config},
};
use std::env;
use std::sync::Arc;
use tokio::time::Duration;
use tracing::{error, info, warn};

mod error;
use error::Error as ReconcilerError;

const REQUEUE_DURATION_FAST: Duration = Duration::from_secs(10);
const REQUEUE_DURATION_SLOW: Duration = Duration::from_secs(60);
const DEFAULT_NAME: &str = "default-name";

mod configmap;
mod crd;
mod deployment;
mod finalizer;
mod helpers;
mod webserver;

#[derive(Clone)]
struct Data {
    client: Client,
    //state: Arc<RwLock<State>>,
}

struct Environment {
    web_server_port: String,
    logger_image: String,
}

struct ImageData {
    image: String,
    version: String,
    policy: String,
}

pub enum NextAction {
    Create,
    Update,
    RecreateConfigMap,
    RecreateDeployment,
    Delete,
    NoAction,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    let kclient = Client::try_default()
        .await
        .expect("Failed to create Kubernetes client");

    // Create an API for the GuaLogger resource
    let loggers: Api<GuaLogger> = Api::all(kclient.clone());

    let context: Arc<Data> = Arc::new(Data {
        client: kclient.clone(),
    });

    tokio::spawn(webserver::create(kclient.clone()));

    Controller::new(loggers.clone(), Config::default())
        .run(reconciler, on_error, context)
        .for_each(|reconciliation_result| async move {
            match reconciliation_result {
                Ok(logger) => {
                    //println!("Reconciliation successful. Resource: {:?}", logger);
                }
                Err(reconciliation_err) => match reconciliation_err {
                    kube::runtime::controller::Error::ObjectNotFound(object_ref) => {
                        warn!("Object not found: {:?}", object_ref);
                    }
                    kube::runtime::controller::Error::ReconcilerFailed(err, object_ref) => {
                        error!(
                            "Reconciler failed for object: {:?}, error: {}",
                            object_ref, err
                        );
                    }
                    kube::runtime::controller::Error::QueueError(err) => {
                        error!("Queue error occurred during reconciliation: {}", err);
                    }
                    kube::runtime::controller::Error::RunnerError(error) => {
                        error!("Runner error: {:?}", error);
                    }
                },
            }
        })
        .await;
}

fn parse_env() -> Environment {
    let port = env::var("WEBSERVER_PORT").unwrap_or("8080".to_owned());
    let image = env::var("LOGGER_IMAGE_SOURCE").unwrap_or("doteich/geist-logger".to_owned());
    let version = env::var("LOGGER_IMAGE_VERSION").unwrap_or("latest".to_owned());
    
}

async fn reconciler(logger: Arc<GuaLogger>, context: Arc<Data>) -> Result<Action, error::Error> {
    let client = &context.client;

    let source_namespace = logger.metadata.namespace.as_deref().unwrap_or("default");

    let name: &str = logger.metadata.name.as_deref().unwrap_or(DEFAULT_NAME);

    let kube_config = extract_kube_struct_from_logger(&logger)?;

    let image = extract_image(kube_config.image);
    let image_str = format!("{}:{}", image.image, image.version);

    let data = extract_data_from_logger(&logger)?;

    match determine_action(&logger, client, &kube_config.namespace).await {
        NextAction::Create => {
            info!("Creating new resource: {:?}", name);
            finalizer::add(client, name, source_namespace).await?;

            configmap::create(client, &kube_config.namespace, name, &data).await?;

            deployment::create(client, &kube_config.namespace, name, &image_str).await?;

            // Here you would typically create the resource, e.g., a deployment
            // For now, we just return a requeue action
            return Ok(Action::requeue(REQUEUE_DURATION_FAST));
        }
        NextAction::Update => {
            info!("Updating existing resource: {:?}", name);
            // Handle update logic here
            return Ok(Action::requeue(REQUEUE_DURATION_FAST));
        }
        NextAction::Delete => {
            finalizer::remove(client, name, source_namespace).await?;

            configmap::delete(client, &kube_config.namespace, name).await?;
            deployment::delete(client, &kube_config.namespace, name).await?;

            // Handle deletion logic here
            return Ok(Action::requeue(REQUEUE_DURATION_FAST));
        }
        NextAction::RecreateConfigMap => {
            info!("Recreating configmap for resource: {:?}", name);

            configmap::create(client, &kube_config.namespace, name, &data).await?;
            return Ok(Action::requeue(REQUEUE_DURATION_FAST));
        }

        NextAction::RecreateDeployment => {
            info!("Recreating deployment for resource: {:?}", name);

            deployment::create(client, &kube_config.namespace, name, &image_str).await?;
            return Ok(Action::requeue(REQUEUE_DURATION_FAST));
        }

        NextAction::NoAction => {
            //info!("No action needed for resource: {:?}", logger.metadata.name);
        }
    }

    Ok(Action::await_change())
}

fn on_error(logger: Arc<GuaLogger>, error: &ReconcilerError, _context: Arc<Data>) -> Action {
    error!(
        "Reconciliation error:\n{:?}.\n{:?}",
        error, logger.metadata.name
    );
    Action::requeue(REQUEUE_DURATION_SLOW)
}

async fn determine_action(logger: &GuaLogger, client: &Client, ns: &str) -> NextAction {
    if logger.metadata.deletion_timestamp.is_some() {
        return NextAction::Delete;
    }

    if logger
        .meta()
        .finalizers
        .as_ref()
        .map_or(true, |finalizers| finalizers.is_empty())
    {
        return NextAction::Create;
    }

    let cm_result = configmap::verify(
        client,
        ns,
        logger.metadata.name.as_deref().unwrap_or("default-name"),
    )
    .await;

    match cm_result {
        Ok(exists) => {
            if !exists {
                return NextAction::RecreateConfigMap;
            }
        }
        Err(e) => {
            error!("Error verifying ConfigMap: {:?}", e);
        }
    }

    let dep_result = deployment::verify(
        client,
        ns,
        logger.metadata.name.as_deref().unwrap_or("default-name"),
    )
    .await;

    match dep_result {
        Ok(exists) => {
            if !exists {
                return NextAction::RecreateDeployment;
            }
        }
        Err(e) => {
            error!("Error verifying ConfigMap: {:?}", e);
        }
    }

    NextAction::NoAction
}

fn extract_data_from_logger(logger: &GuaLogger) -> Result<String, ReconcilerError> {
    if let Some(logger_data) = &logger.spec.logger {
        let yaml_raw = serde_yaml::to_string(logger_data)?;
        Ok(yaml_raw)
    } else {
        Err(ReconcilerError::MissingLoggerData)
    }
}

fn extract_kube_struct_from_logger(
    logger: &GuaLogger,
) -> Result<crd::GuaLoggerKube, ReconcilerError> {
    if let Some(kube_struct) = &logger.spec.kube {
        Ok(kube_struct.clone())
    } else {
        Err(ReconcilerError::MissingKubeObject)
    }
}

fn extract_image(kube_str: Option<GuaLoggerKubeImage>) -> ImageData {
    let mut data = {
        ImageData {
            image: "cinderstries/gualogger".to_string(),
            version: "0.0.1".to_string(),
            policy: "Always".to_string(),
        }
    };

    match kube_str {
        Some(k) => {
            if let Some(repo) = k.repository {
                data.image = repo;
            };
            if let Some(version) = k.version {
                data.version = version;
            };

            if let Some(policy) = k.pull_policy {
                data.policy = policy
            }

            data
        }
        None => data,
    }
}
