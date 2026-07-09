pipeline {
    agent any
    
    stages {
        stage('Grab Latest Code') {
            steps {
                
                echo 'Pulling the latest updates from GitHub...'
                checkout scm
            }
        }
        
        stage('Build Docker Image') {
            steps {
               
                echo 'Starting the multi-stage docker build for the app...'
                sh 'docker build -t inventory-app:latest .'
            }
        }
        
        stage('Fresh Deployment') {
            steps {
                
                echo 'Clearing out old instances and spinning up the updated container...'
                sh 'docker rm -f inventory-app-container || true'
                sh 'docker run -d --name inventory-app-container -p 8080:8080 inventory-app:latest'
            }
        }
    }
}
