FROM python:3.9-slim

WORKDIR /app

# Copy the AI service placeholder
COPY ./ai-placeholder/placeholder_server.py /app/

# Install required packages
RUN pip install --no-cache-dir flask

# Expose port
EXPOSE 9000

# Set environment variables
ENV PORT=9000
ENV MODEL_PATH=/app/models
ENV DATA_PATH=/app/data

# Run the server
CMD ["python", "/app/placeholder_server.py"]
