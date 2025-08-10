use axum::{
    Json, Router,
    extract::State,
    http,
    response::{IntoResponse, IntoResponseParts},
    routing::get,
};
use k8s_openapi::api::apps::v1::Deployment;
use kube::{api::ObjectList, client};
use std::error::Error;

use crate::{
    crd::{self, GuaLogger},
    deployment,
};

struct CustomError {
    message: String,
    status_code: http::StatusCode,
}
impl IntoResponse for CustomError {
    fn into_response(self) -> axum::response::Response {
        (
            self.status_code,
            axum::Json(serde_json::json!({ "error": self.message })),
        )
            .into_response()
    }
}

pub async fn create(client: kube::Client) {
    let router: Router = Router::new()
        .route(
            "/api/crds",
            get(fetch_crds),
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

async fn fetch_deployments(
    state: State<kube::Client>,
) -> Result<Json<ObjectList<Deployment>>, CustomError> {
    match deployment::get(&state.clone()).await {
        Ok(res) => Ok(Json(res)),
        Err(e) => {
            println!("{}", e);
            let err = CustomError {
                message: "error while retrieving deployments".to_string(),
                status_code: http::StatusCode::INTERNAL_SERVER_ERROR,
            };
            Err(err)
        }
    }
}

async fn fetch_crds(
    state: State<kube::Client>,
) -> Result<Json<ObjectList<GuaLogger>>, CustomError> {
    match crd::get(&state.clone()).await {
        Ok(res) => Ok(Json(res)),
        Err(e) => {
            println!("{}", e);
            let err = CustomError {
                message: "error while retrieving deployments".to_string(),
                status_code: http::StatusCode::INTERNAL_SERVER_ERROR,
            };
            Err(err)
        }
    }
}
