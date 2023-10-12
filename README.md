
# Upload file to S3 Bucket

This project provides a command-line tool for uploading files to an Amazon Web Services (AWS) S3 bucket and reading the uploaded files.

##  Features
- Upload files to an S3 bucket
- Support for specifying a target folder within the bucket
- Read uploaded files from an S3 bucket

## Getting Started

### Prerequisites
- An AWS account with access to an S3 bucket
- AWS credentials with permissions to access the target bucket

### Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/ishangoyal13/aws-uploader
```

2. Install the dependencies:   
   
```bash
cd aws-uploader
go mod download
```

3. Run the project:  

```go
go run main.go
```

