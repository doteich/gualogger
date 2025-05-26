use std::sync::{Arc, RwLock};

use crd::GuaLogger;
use k8s_openapi::api::core::v1::Pod;
use kube::Client;
use kube::Error;
use kube::api::{Api, ListParams};
use kube::runtime::Controller;
use kube::runtime::controller::Action;
use kube::runtime::reflector::Lookup;
use kube::runtime::watcher::Config;
use serde::de;
use std::error::Error as StdError;
use std::future::Future;

use futures::stream::StreamExt;
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

    println!("got update for object: {:?}, with uid: {:?}", logger.metadata.name, logger.metadata.uid);

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

fn on_error(echo: Arc<GuaLogger>, error: &Error, _context: Arc<Data>) -> Action {
    eprintln!("Reconciliation error:\n{:?}.\n{:?}", error, echo);
    Action::requeue(Duration::from_secs(5))
}
