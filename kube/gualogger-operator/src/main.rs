use std::sync::{Arc, RwLock};

use crd::GuaLogger;
use k8s_openapi::api::core::v1::Pod;
use kube::api::{Api, ListParams};
use kube::runtime::controller::Action;
use kube::runtime::reflector::Lookup;
use kube::runtime::watcher::Config;
use kube::runtime::Controller;
use kube::Client;
use kube::Error;
use serde::de;
use std::future::Future;
use std::error::Error as StdError;

use tokio::time::Duration;
use futures::stream::StreamExt;

mod crd;
mod deployment;
mod helpers;

#[derive(Clone)]
struct Data {
    client: Client,
    //state: Arc<RwLock<State>>,
}

#[tokio::main]
async fn main() {
    // Load the kubeconfig file from the default location
    // let kubeconfig_path = env::var("KUBECONFIG").unwrap_or_else(|_| {
    //     let home_dir = env::var("HOME").unwrap();
    //     format!("{}/.kube/config", home_dir)
    // });

    // Create a Kubernetes client
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
                    println!("Reconciliation successful. Resource: {:?}", logger);
                }
                Err(reconciliation_err) => {
                    eprintln!("Reconciliation error: {:?}", reconciliation_err)
                }
            }
        })
        .await;
}

async fn reconciler(
    logger: Arc<GuaLogger>,
    context: Arc<Data>,
) -> Result<Action, Error> {
    // Implement your reconciliation logic here
    // For example, you can use the context to access the Kubernetes client
    let client = &context.client;

    // Perform some operations with the client
    // ...

    let res = deployment::spawn_deployment(
        client,
        "default",
        logger.name().as_deref().unwrap_or("default-name"),
        "nginx:latest",
    )
    .await;

    match res {
        Ok(_) => {
            Ok(Action::requeue(Duration::from_secs(10)))
        }
        Err(e) => {
            eprintln!("Error creating deployment: {:?}", e);
           Err(e)
        }
    }
}

fn on_error(echo: Arc<GuaLogger>, error: &Error, _context: Arc<Data>) -> Action {
    eprintln!("Reconciliation error:\n{:?}.\n{:?}", error, echo);
    Action::requeue(Duration::from_secs(5))
}