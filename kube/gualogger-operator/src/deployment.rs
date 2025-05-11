use std::collections::BTreeMap;

use k8s_openapi::api::apps::v1::{Deployment, DeploymentSpec};
use k8s_openapi::api::core::v1::{Container, ContainerPort, PodSpec, PodTemplateSpec};
use k8s_openapi::apimachinery::pkg::apis::meta::v1::LabelSelector;
use kube::{Api, Client};
use kube::api::{ObjectMeta, PostParams};

pub async fn spawn_deployment(
    client: &Client,
    namespace: &str,
    name: &str,
    image: &str,
) -> Result<(), kube::Error> {


    let mut labels = Some(BTreeMap::new());


    labels.as_mut().unwrap().insert("app".to_string(), name.to_string());

    // Create a new deployment
    let deployment = Deployment {
        metadata: ObjectMeta {
            name: Some(name.to_string()),
            namespace: Some(namespace.to_string()),
            labels: labels.clone(),
            ..ObjectMeta::default()
           
        },
        spec: Some(DeploymentSpec {
            replicas: Some(1),
            selector: LabelSelector {
                match_labels: labels.clone(),
                ..LabelSelector::default()
            },
            template: PodTemplateSpec {
                metadata: Some(ObjectMeta {
                    labels: labels.clone(),
                    ..ObjectMeta::default()
                }),
                spec: Some(PodSpec {
                    containers: vec![Container {
                        name: name.to_string(),
                        image: Some(image.to_string()),
                        ..Container::default()
                    }],
                    ..PodSpec::default()
                }),
            },
            ..DeploymentSpec::default()
        }),
        ..Deployment::default()
    };

    // Create the deployment in the specified namespace
    let deployments = Api::<Deployment>::namespaced(client.clone(), namespace);
    deployments
        .create(&PostParams::default(), &deployment)
        .await?;

    Ok(())
}
