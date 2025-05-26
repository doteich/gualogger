use kube::{Api, Client};

pub async fn add(client: Client) -> Result<(), kube::Error> {
    // ToDo:
    // Fetch specific CRD GuaLogger
    // Add finalizer to the GuaLogger CRD by patching it
    // See https://github.com/Pscheidl/rust-kubernetes-operator-example/blob/master/src/finalizer.rs

    Ok(())
}
