use std::sync::{Arc, RwLock};

use crd::GuaLogger;

use kube::Client;
use kube::Error;
use kube::Resource;
use kube::ResourceExt;
use kube::api::Api;
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

    println!(
        "got update for object: {:?}, with uid: {:?}",
        logger.metadata.name, logger.metadata.uid
    );

    let namespace = logger.metadata.namespace.as_deref().unwrap_or("default");

    match determine_action(&logger) {
        NextAction::Create => {
            println!("Creating new resource: {:?}", logger.metadata.name);
            finalizer::add(
                client.clone(),
                logger.metadata.name.as_deref().unwrap_or("default-name"),
                namespace,
            )
            .await?;
            // Here you would typically create the resource, e.g., a deployment
            // For now, we just return a requeue action
            return Ok(Action::requeue(Duration::from_secs(10)));
        }
        NextAction::Update => {
            println!("Updating existing resource: {:?}", logger.metadata.name);
            // Handle update logic here
            return Ok(Action::requeue(Duration::from_secs(10)));
        }
        NextAction::Delete => {
            println!("Deleting resource: {:?}", logger.metadata.name);
            finalizer::remove(
                client.clone(),
                logger.metadata.name.as_deref().unwrap_or("default-name"),
                namespace,
            )
            .await?;
            // Handle deletion logic here
            return Ok(Action::requeue(Duration::from_secs(10)));
        }
        NextAction::NoAction => {
            println!("No action needed for resource: {:?}", logger.metadata.name);
        }
    }

    // let res = deployment::spawn_deployment(
    //     client,
    //     "default",
    //     logger.name().as_deref().unwrap_or("default-name"),
    //     "nginx:latest",
    // )
    // .await;

    // match res {
    //     Ok(_) => Ok(Action::requeue(Duration::from_secs(10))),
    //     Err(e) => {
    //         eprintln!("Error creating deployment: {:?}", e);
    //         Err(e)
    //     }
    // }

    Ok(Action::requeue(Duration::from_secs(10)))
}

fn on_error(logger: Arc<GuaLogger>, error: &Error, _context: Arc<Data>) -> Action {
    eprintln!(
        "Reconciliation error:\n{:?}.\n{:?}",
        error, logger.metadata.name
    );
    Action::requeue(Duration::from_secs(60))
}

fn determine_action(logger: &GuaLogger) -> NextAction {
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

    NextAction::NoAction
}
