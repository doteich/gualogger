use crate::crd::GuaLogger;
use kube::{
    Api, Client,
    api::{Patch, PatchParams},
    runtime::reflector::Lookup,
};
use serde_json::{Value, json};

pub async fn add(client: Client, name: &str, namespace: &str) -> Result<(), kube::Error> {
    // ToDo:
    // Fetch specific CRD GuaLogger
    // Add finalizer to the GuaLogger CRD by patching it
    // See https://github.com/Pscheidl/rust-kubernetes-operator-example/blob/master/src/finalizer.rs

    let loggers: Api<GuaLogger> = Api::namespaced(client, namespace);

    let finalizer: Value = json!({
        "metadata": {
            "finalizers": ["gualoggers.doteich.com/finalizer"]
        }
    });

    let patch: Patch<&Value> = Patch::Merge(&finalizer);

    println!("Fetching logger: {}", name);

    let logger = loggers.get(name).await?;

    println!("Adding finalizer to logger: {}", name);

    if let Some(name) = logger.name() {
        loggers
            .patch(&name, &PatchParams::default(), &patch)
            .await?;
    } else {
        println!("Logger {} not found", name);
    }

    Ok(())
}
