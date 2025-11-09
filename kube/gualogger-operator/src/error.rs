use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("Kubernetes error: {0}")]
    KubeError(#[from] kube::Error),

    #[error("Serde yaml error: {0}")]
    SerdeYamlError(#[from] serde_yaml::Error),

    #[error("Missing kube object in spec")]
    MissingKubeObject,

    #[error("Missing logger data in spec")]
    MissingLoggerData,
}
