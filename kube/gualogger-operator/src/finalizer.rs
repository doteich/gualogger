use crate::crd::GuaLogger;
use futures::future::err;
use kube::{
    Api, Client,
    api::{Patch, PatchParams},
    runtime::reflector::Lookup,
};

use kube::error::ErrorResponse;
use serde_json::{Value, json};

pub async fn add(client: &Client, name: &str, namespace: &str) -> Result<(), kube::Error> {
    // ToDo:
    // Fetch specific CRD GuaLogger
    // Add finalizer to the GuaLogger CRD by patching it
    // See https://github.com/Pscheidl/rust-kubernetes-operator-example/blob/master/src/finalizer.rs

    let loggers: Api<GuaLogger> = Api::namespaced(client.clone(), namespace);

    let finalizer: Value = json!({
        "metadata": {
            "finalizers": ["gualoggers.doteich.com/finalizer"]
        }
    });

    let patch: Patch<&Value> = Patch::Merge(&finalizer);

    let logger = loggers.get(name).await?;

    if let Some(name) = logger.name() {
        loggers
            .patch(&name, &PatchParams::default(), &patch)
            .await?;
    } else {
        return Err(kube::Error::Api(ErrorResponse {
            status: "logger not found".to_string(),
            message: "could not find a suitable logger with provided parameters".to_string(),
            reason: format!("logger {} not found in namespace {}", name, namespace),
            code: 404,
        }));
    }

    Ok(())
}

pub async fn remove(client: &Client, name: &str, namespace: &str) -> Result<(), kube::Error> {
    let loggers: Api<GuaLogger> = Api::namespaced(client.clone(), namespace);

    let finalizer: Value = json!({
        "metadata": {
            "finalizers": []
        }
    });

        let patch: Patch<&Value> = Patch::Merge(&finalizer);

    let logger = loggers.get(name).await?;

    if let Some(name) = logger.name() {
        loggers
            .patch(&name, &PatchParams::default(), &patch)
            .await?;
    } else {
        return Err(kube::Error::Api(ErrorResponse {
            status: "logger not found".to_string(),
            message: "could not find a suitable logger with provided parameters".to_string(),
            reason: format!("logger {} not found in namespace {}", name, namespace),
            code: 404,
        }));
    }

    Ok(())
}
