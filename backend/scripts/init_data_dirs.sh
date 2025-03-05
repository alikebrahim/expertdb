#!/bin/bash

# This script initializes the data directory structure for ExpertDB

# Create base data directory
mkdir -p data

# Create SQLite database directory
mkdir -p data/sqlite

# Create document storage directories
mkdir -p data/documents/cv
mkdir -p data/documents/certificate
mkdir -p data/documents/publication

# Set permissions
chmod -R 755 data

echo "Data directory structure initialized successfully."
echo "Directory structure:"
find data -type d | sort