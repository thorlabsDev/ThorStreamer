// lib.rs
use tonic::transport::{Channel, Endpoint};
use tonic::{metadata::MetadataValue, Request, Status, Streaming};
use std::time::Duration;

// Generated protobuf modules
pub mod proto {
    pub mod publisher {
        tonic::include_proto!("publisher");
    }
    pub mod thor_streamer {
        pub mod types {
            tonic::include_proto!("thor_streamer.types");
        }
    }
    pub mod google {
        pub mod protobuf {
            tonic::include_proto!("google.protobuf");
        }
    }
}

// Import both Empty types with aliases
use proto::google::protobuf::Empty as GoogleEmpty;
use proto::thor_streamer::types::Empty as ThorEmpty;

use proto::publisher::{
    event_publisher_client::EventPublisherClient,
    SubscribeAccountsRequest, StreamResponse,
};
use proto::thor_streamer::types::{
    thor_streamer_client::ThorStreamerClient,
    MessageWrapper,
};

/// Configuration for the Thor gRPC client
#[derive(Clone, Debug)]
pub struct ClientConfig {
    pub server_addr: String,
    pub token: String,
    pub timeout: Duration,
}

impl Default for ClientConfig {
    fn default() -> Self {
        Self {
            server_addr: "http://localhost:50051".to_string(),
            token: String::new(),
            timeout: Duration::from_secs(30),
        }
    }
}

/// Main client for Thor gRPC services
#[derive(Clone)]
pub struct ThorClient {
    event_client: EventPublisherClient<Channel>,
    thor_client: ThorStreamerClient<Channel>,
    token: String,
}

impl ThorClient {
    /// Create a new Thor client
    pub async fn new(config: ClientConfig) -> Result<Self, Box<dyn std::error::Error>> {
        let endpoint = Endpoint::from_shared(config.server_addr)?
            .timeout(config.timeout)
            .tcp_keepalive(Some(Duration::from_secs(30)))
            .http2_keep_alive_interval(Duration::from_secs(30))
            .keep_alive_timeout(Duration::from_secs(10));

        let channel = endpoint.connect().await?;

        Ok(Self {
            event_client: EventPublisherClient::new(channel.clone()),
            thor_client: ThorStreamerClient::new(channel),
            token: config.token,
        })
    }

    /// Create a request with authentication
    fn with_auth<T>(&self, req: T) -> Request<T> {
        let mut request = Request::new(req);
        if let Ok(token) = MetadataValue::try_from(&self.token) {
            request.metadata_mut().insert("authorization", token);
        }
        request
    }

    /// Subscribe to transaction events (uses GoogleEmpty)
    pub async fn subscribe_to_transactions(&mut self) -> Result<Streaming<StreamResponse>, Status> {
        let request = self.with_auth(GoogleEmpty {});
        let response = self.event_client.subscribe_to_transactions(request).await?;
        Ok(response.into_inner())
    }

    /// Subscribe to slot status events (uses GoogleEmpty)
    pub async fn subscribe_to_slot_status(&mut self) -> Result<Streaming<StreamResponse>, Status> {
        let request = self.with_auth(GoogleEmpty {});
        let response = self.event_client.subscribe_to_slot_status(request).await?;
        Ok(response.into_inner())
    }

    /// Subscribe to account update events
    pub async fn subscribe_to_account_updates(
        &mut self,
        account_addresses: Vec<String>,
        owner_addresses: Vec<String>,
    ) -> Result<Streaming<StreamResponse>, Status> {
        let request = self.with_auth(SubscribeAccountsRequest {
            account_address: account_addresses,
            owner_address: owner_addresses,
        });
        let response = self
            .event_client
            .subscribe_to_account_updates(request)
            .await?;
        Ok(response.into_inner())
    }

    /// Subscribe to Thor update stream (uses ThorEmpty)
    pub async fn subscribe_to_thor_updates(&mut self) -> Result<Streaming<MessageWrapper>, Status> {
        let request = Request::new(ThorEmpty {});
        let response = self.thor_client.stream_updates(request).await?;
        Ok(response.into_inner())
    }
}

/// Parse a StreamResponse into MessageWrapper
pub fn parse_message(data: &[u8]) -> Result<MessageWrapper, prost::DecodeError> {
    use prost::Message;
    MessageWrapper::decode(data)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_config_default() {
        let config = ClientConfig::default();
        assert_eq!(config.server_addr, "http://localhost:50051");
        assert_eq!(config.timeout, Duration::from_secs(30));
    }
}