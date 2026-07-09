pipeline {
    agent any
    
    stages {
        stage('SCM Checkout') {
            steps {
                echo 'Pulling latest code from GitHub...'
                checkout scm
            }
        }
        
        stage('Docker Build') {
            steps {
                echo 'Building Multi-Stage Docker Image...'
                sh 'docker compose build'
            }
        }
        
        stage('Infrastructure Deploy') {
            steps {
                echo 'Deploying Go API and Postgres Database Containers...'
                sh 'docker compose down'
                sh 'docker compose up -d'
            }
        }
    }
}
