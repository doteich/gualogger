use crate::crd::GuaLogger;
use k8s_openapi::api::core::v1::ConfigMap;
use kube::api::{ObjectMeta, PostParams};
use kube::Client;
use std::collections::BTreeMap;

fn create_configmap(client: &Client, namespace: &str, name: &str, data: &String) {
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
}
