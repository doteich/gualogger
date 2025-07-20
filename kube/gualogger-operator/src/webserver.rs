use axum::{Json, Router, extract::State, routing::get};
use k8s_openapi::api::apps::v1::Deployment;
use kube::{api::ObjectList, client};
use std::error::Error;

use crate::deployment;

pub async fn create(client: kube::Client) {
    let router: Router = Router::new()
        .route(
            "/api/crds",
            get(|| async {
                println!("Received inbound request");
                "Hello, World!"
            }),
        )
        .route("/api/resources", get(fetch_deployments))
        .with_state(client);

    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
        .await
        .unwrap();
    println!("server created");

    match axum::serve(listener, router).await {
        Ok(_) => (),
        Err(err) => {
            println!("{}", err)
        }
    }
}

async fn fetch_deployments(state: State<kube::Client>) -> Json<ObjectList<Deployment>> {
    let res = deployment::get(&&state.clone()).await;
    let j = Json(res);
    return j;
}
