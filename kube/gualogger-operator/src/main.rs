use std::sync::{Arc, RwLock};

use crd::GuaLogger;

use kube::Client;
use kube::Error;
use kube::Resource;
use kube::ResourceExt;
use kube::api::Api;
use kube::client;
use kube::config;
use kube::runtime::Controller;
use kube::runtime::controller::Action;
use kube::runtime::reflector::Lookup;
use kube::runtime::watcher::Config;

use futures::stream::StreamExt;
use serde::de;
use tokio::time::Duration;

mod configmap;
mod crd;
mod deployment;
mod finalizer;
mod helpers;

#[derive(Clone)]
struct Data {
    client: Client,
    //state: Arc<RwLock<State>>,
}

pub enum NextAction {
    Create,
    Update,
    RecreateConfigMap,
    // RecreateDeployment,
    Delete,
    NoAction,
}

#[tokio::main]
async fn main() {
    let kclient = Client::try_default()
        .await
        .expect("Failed to create Kubernetes client");

    // Create an API for the GuaLogger resource
    let loggers: Api<GuaLogger> = Api::all(kclient.clone());

    let context: Arc<Data> = Arc::new(Data {
        client: kclient.clone(),
    });

    Controller::new(loggers.clone(), Config::default())
        .run(reconciler, on_error, context)
        .for_each(|reconciliation_result| async move {
            match reconciliation_result {
                Ok(logger) => {
                    //println!("Reconciliation successful. Resource: {:?}", logger);
                }
                Err(reconciliation_err) => match reconciliation_err {
                    kube::runtime::controller::Error::ObjectNotFound(object_ref) => {
                        println!("Object not found: {:?}", object_ref);
                    }
                    kube::runtime::controller::Error::ReconcilerFailed(_, object_ref) => {
                        println!("Reconciler failed for object: {:?}", object_ref);
                    }
                    kube::runtime::controller::Error::QueueError(_) => {
                        eprintln!("Queue error occurred during reconciliation.");
                    }
                    kube::runtime::controller::Error::RunnerError(error) => {
                        eprintln!("Runner error: {:?}", error);
                    }
                },
            }
        })
        .await;
}

async fn reconciler(logger: Arc<GuaLogger>, context: Arc<Data>) -> Result<Action, Error> {
    let client = &context.client;

    let source_namespace = logger.metadata.namespace.as_deref().unwrap_or("default");

    let name: &str = logger.metadata.name.as_deref().unwrap_or("default-name");

    let kube_config = extract_kube_struct_from_logger(&logger.clone()).map_err(|e| {
        Error::from(kube::Error::Api(kube::core::ErrorResponse {
            status: "bad request".to_string(),
            message: e.to_string(),
            reason: "kube object is missing from spec".to_string(),
            code: 404,
        }))
    })?;
    let data = extract_data_from_logger(&logger).map_err(|e| {
        Error::from(kube::Error::Api(kube::core::ErrorResponse {
            status: "bad request".to_string(),
            message: e.to_string(),
            reason: "logger data is missing from spec".to_string(),
            code: 404,
        }))
    })?;

    match determine_action(&logger, client, &kube_config.namespace).await {
        NextAction::Create => {
            println!("Creating new resource: {:?}", name);
            finalizer::add(client.clone(), name, source_namespace).await?;

            configmap::create(
                client,
                &kube_config.namespace,
                format!("{}-configmap", name).as_str(),
                &data,
            )
            .await?;

            // Here you would typically create the resource, e.g., a deployment
            // For now, we just return a requeue action
            return Ok(Action::requeue(Duration::from_secs(10)));
        }
        NextAction::Update => {
            println!("Updating existing resource: {:?}", name);
            // Handle update logic here
            return Ok(Action::requeue(Duration::from_secs(10)));
        }
        NextAction::Delete => {
            println!("Deleting resource: {:?}", name);
            finalizer::remove(client.clone(), name, source_namespace).await?;
            // Handle deletion logic here
            return Ok(Action::requeue(Duration::from_secs(10)));
        }
        NextAction::RecreateConfigMap => {
            println!("Recreating ConfigMap for resource: {:?}", name);

            configmap::create(
                client,
                &kube_config.namespace,
                format!("{}-configmap", name).as_str(),
                &data,
            )
            .await?;
            return Ok(Action::requeue(Duration::from_secs(10)));
        }
        NextAction::NoAction => {
            println!("No action needed for resource: {:?}", logger.metadata.name);
        }
    }

    Ok(Action::requeue(Duration::from_secs(10)))
}

fn on_error(logger: Arc<GuaLogger>, error: &Error, _context: Arc<Data>) -> Action {
    eprintln!(
        "Reconciliation error:\n{:?}.\n{:?}",
        error, logger.metadata.name
    );
    Action::requeue(Duration::from_secs(60))
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

    let result = configmap::verify(
        client,
        ns,
        &format!(
            "{}-configmap",
            logger.metadata.name.as_deref().unwrap_or("default-name")
        ),
    )
    .await;

    match result {
        Ok(exists) => {
            if !exists {
             
                return NextAction::RecreateConfigMap;
            }
        }
        Err(e) => {
            eprintln!("Error verifying ConfigMap: {:?}", e);
        }
    }

    NextAction::NoAction
}

fn extract_data_from_logger(logger: &GuaLogger) -> Result<String, Box<dyn std::error::Error>> {
    if let Some(logger_data) = &logger.spec.logger {
        let yaml_raw = serde_yaml::to_string(logger_data)?;
        Ok(yaml_raw)
    } else {
        Err("No logger data found in logger spec".into())
    }
}

fn extract_kube_struct_from_logger(
    logger: &GuaLogger,
) -> Result<crd::GuaLoggerKube, Box<dyn std::error::Error>> {
    if let Some(kube_struct) = &logger.spec.kube {
        Ok(kube_struct.clone())
    } else {
        Err("No kube struct found in logger spec".into())
    }
}
