use std::collections::BTreeMap;

use k8s_openapi::api::apps::v1::{Deployment, DeploymentSpec};
use k8s_openapi::api::core::v1::{
    ConfigMapVolumeSource, Container, PodSpec, PodTemplateSpec, Volume, VolumeMount,
};
use k8s_openapi::apimachinery::pkg::apis::meta::v1::LabelSelector;
use kube::api::{DeleteParams, ObjectMeta, PostParams};
use kube::{Api, Client};

pub async fn create(
    client: &Client,
    namespace: &str,
    name: &str,
    image: &str,
) -> Result<(), kube::Error> {
    let mut labels = Some(BTreeMap::new());

    labels
        .as_mut()
        .unwrap()
        .insert("app".to_string(), name.to_string());

    // Create a new deployment
    let deployment = Deployment {
        metadata: ObjectMeta {
            name: Some(name.to_string() + "-deployment"),
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
                        volume_mounts: Some(vec![VolumeMount {
                            mount_path: "/etc/gopclogs".to_string(),
                            name: "gopclogs-volume".to_string(),
                            ..Default::default()
                        }]),

                        ..Default::default()
                    }],

                    volumes: Some(vec![Volume {
                        name: "gopclogs-volume".to_string(),
                        config_map: Some(ConfigMapVolumeSource {
                            name: name.to_string() + "-configmap",
                            ..Default::default()
                        }),
                        ..Default::default()
                    }]),
                    ..Default::default()
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

pub async fn verify(client: &Client, namespace: &str, name: &str) -> Result<bool, kube::Error> {
    let configmaps: kube::Api<Deployment> = kube::Api::namespaced(client.clone(), namespace);

    let cm_name = name.to_owned() + "-deployment";

    match configmaps.get(&cm_name).await {
        Ok(_) => Ok(true),
        Err(kube::Error::Api(e)) if e.code == 404 => Ok(false),
        Err(e) => Err(e),
    }
}

pub async fn delete(client: &Client, namespace: &str, name: &str) -> Result<(), kube::Error> {
    let configmaps: kube::Api<Deployment> = kube::Api::namespaced(client.clone(), namespace);

    let d_name = name.to_owned() + "-deployment";

    if let Err(e) = configmaps.delete(&d_name, &DeleteParams::default()).await {
        return Err(e);
    } else {
        Ok(())
    }
}
