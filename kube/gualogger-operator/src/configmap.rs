use crate::crd::GuaLogger;
use k8s_openapi::api::core::v1::ConfigMap;
use kube::Client;
use kube::api::{ObjectMeta, PostParams};
use std::collections::BTreeMap;

pub async fn create(client: &Client, namespace: &str, name: &str, data: &String) -> Result<(), kube::Error> {
    let mut labels = Some(BTreeMap::new());

    labels
        .as_mut()
        .unwrap()
        .insert("app".to_string(), name.to_string());

    let configmap = ConfigMap {
        metadata: ObjectMeta {
            name: Some(name.to_string()),
            namespace: Some(namespace.to_string()),
            ..ObjectMeta::default()
        },
        data: Some(BTreeMap::from([("config.yaml".to_string(), data.clone())])),
        ..ConfigMap::default()
    };

    let configmaps: kube::Api<ConfigMap> = kube::Api::namespaced(client.clone(), namespace);

    configmaps.create(&PostParams::default(), &configmap).await?;

    Ok(())
}


pub async fn verify(client: &Client, namespace: &str, name: &str) -> Result<bool, kube::Error> {
    let configmaps: kube::Api<ConfigMap> = kube::Api::namespaced(client.clone(), namespace);
    match configmaps.get(name).await {
        Ok(_) => Ok(true),
        Err(kube::Error::Api(e)) if e.code == 404 => Ok(false),
        Err(e) => Err(e),
    }
}

// pub async fn delete(client: &Client, namespace: &str, name: &str) {
//     let configmaps: kube::Api<ConfigMap> = kube::Api::namespaced(client.clone(), namespace);
//     if let Err(e) = configmaps.delete(name, &PostParams::default()).await {
//         eprintln!("Failed to delete ConfigMap {}: {}", name, e);
//     }
// }
