use kube::CustomResourceExt;
use super::crd;


fn main() {
    print!("{}", serde_yaml::to_string(&crd::GuaLogger::crd()).unwrap())
}