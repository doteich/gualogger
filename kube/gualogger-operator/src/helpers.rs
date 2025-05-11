use crate::crd;
use kube::CustomResourceExt;

pub fn create_crd() {
    print!("{}", serde_yaml::to_string(&crd::GuaLogger::crd()).unwrap())
}

pub fn align_crd() {
    let crd = crd::GuaLogger::crd();
    let crd_yaml = serde_yaml::to_string(&crd).unwrap();
    let aligned_crd = crd_yaml.replace("  ", "    ");
    print!("{}", aligned_crd);
}
