use k8s_openapi::api::core::{v1::Pod};
use kube::Client;
use kube::api::{Api, ListParams};

use std::env;

#[tokio::main]
async fn main() {
    // Load the kubeconfig file from the default location
    // let kubeconfig_path = env::var("KUBECONFIG").unwrap_or_else(|_| {
    //     let home_dir = env::var("HOME").unwrap();
    //     format!("{}/.kube/config", home_dir)
    // });

    // Create a Kubernetes client
    let client = Client::try_default()
        .await
        .expect("Failed to create Kubernetes client");

    // Create an API for the Pod resource
    let pods: Api<Pod> = Api::all(client);

    // List all Pods in the default namespace
    let lp = ListParams::default();
    match pods.list(&lp).await {
        Ok(pod_list) => {
            for pod in pod_list.items {
                println!("Found Pod: {}", pod.metadata.name.unwrap());
            }
        }
        Err(e) => {
            eprintln!("Error listing Pods: {}", e);
        }
    }
}
