# Leader Election and Mutual Exclusion (Distributed File Replication)

A Go-based distributed system demonstrating the **Bully Algorithm** for leader election, **Ricart-Agrawala Algorithm** for mutual exclusion, and **Consensus-based File Replication**.

## 🚀 Quick Start (Local Demo)

To run a 3-node cluster on a single machine:

1.  **Terminal 1 (Node 1)**:
    ```bash
    go run ./distributed-file-replication -id 1 -port 8001 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003"
    ```
2.  **Terminal 2 (Node 2)**:
    ```bash
    go run ./distributed-file-replication -id 2 -port 8002 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003"
    ```
3.  **Terminal 3 (Node 3)**:
    ```bash
    go run ./distributed-file-replication -id 3 -port 8003 -peers "1=localhost:8001,2=localhost:8002,3=localhost:8003"
    ```

## 💻 Multi-Laptop Demo Setup

To run across 5 laptops:

1.  **Connect to same Wi-Fi** (Mobile Hotspot recommended).
2.  **Find IP of each laptop** (`ipconfig` in CMD).
3.  **Run on each laptop** using its ID and the full peers list:
    ```bash
    # Example for Laptop 1
    go run ./distributed-file-replication -id 1 -port 8001 -peers "1=IP1:8001,2=IP2:8002,3=IP3:8003,4=IP4:8004,5=IP5:8005"
    ```

## 🛠 Interactive Commands

Once a node is running, type these into the terminal:

-   `status`: View Current Node ID, Leader, and Files.
-   `election`: Trigger a new leader election.
-   `cs`: Request access to the Critical Section.
-   `replicate`: Replicate `report.txt` across all active nodes.
-   `exitcs`: Release the Critical Section.

## 📁 System Architecture

-   **Leader Election**: Implemented using the Bully Algorithm.
-   **Mutual Exclusion**: Implemented using the Ricart-Agrawala algorithm (distributed lamport clocks).
-   **Consensus**: Files are only committed if a majority of nodes acknowledge the write.
-   **Storage**: Each node stores its files in the `node_storage/nodeX/` directory.

## ⚠️ Troubleshooting

-   **Firewall**: If nodes can't see each other, allow `go.exe` through Windows Firewall.
-   **Quorum**: In a 5-node setup, at least 3 nodes must be online for replication to succeed.
