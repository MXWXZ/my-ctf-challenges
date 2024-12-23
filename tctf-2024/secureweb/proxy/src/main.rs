use std::env;
use std::net::SocketAddr;

use regex::Regex;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};
use tokio::spawn;

async fn handle_connection(
    mut inbound: TcpStream,
    target_addr: &str,
    addr: SocketAddr,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut outbound = TcpStream::connect(target_addr).await?;

    let mut buffer = vec![0; 4096];
    let bytes_read = inbound.read(&mut buffer).await?;
    let request_data = String::from_utf8_lossy(&buffer[..bytes_read]);

    let mut lines = request_data.lines();

    if let Some(line) = lines.next() {
        let mut req = format!("{}\n", line);
        req.push_str(&format!("X-Forwarded-For: {}\n", addr.ip()));

        for line in lines {
            req.push_str(line);
            req.push_str("\n");
        }

        let re = Regex::new(r"\r|\n").unwrap();
        req = re.replace_all(&req, "\r\n").to_string();

        outbound.write_all(req.as_bytes()).await?;
    }

    let r = outbound.read(&mut buffer).await?;
    inbound.write_all(&buffer[..r]).await?;
    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let listener = TcpListener::bind("0.0.0.0:8080").await?;
    println!("Listening on 0.0.0.0:8080");

    loop {
        let (inbound, addr) = listener.accept().await?;
        spawn(async move {
            if let Err(e) = handle_connection(inbound, &env::var("WEB").unwrap(), addr).await {
                println!("Error handling connection: {e}");
            }
        });
    }
}
