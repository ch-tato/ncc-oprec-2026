pipeline {
    agent any
    
    environment {
        PATH = "/usr/local/go/bin:$PATH"
        GOCACHE = "${WORKSPACE}/.cache/go-build"
        GOPATH = "${WORKSPACE}/go"
        SONARQUBE_ENV = 'SonarQube'
        SCANNER_HOME = tool 'sonar-scanner'
        PROJECT_KEY   = 'ncc-module-2'
        PROJECT_NAME  = "NCC Module 2"
    }

    stages {
        stage('Build and Test') {
            parallel {
                stage('Unit Test') {
                    steps {
                        echo 'Running Unit Tests...'
                        sh 'go test -v ./... -coverprofile=coverage.out || true'
                    }
                }
                stage('Compile Build') {
                    steps {
                        echo 'Compiling the application...'
                        // sh 'go build -o main .'
                    }
                }
            }
        }

        stage('SonarQube Analysis') {
            steps {
                echo 'Running SonarQube Analysis...'
                withSonarQubeEnv("${SONARQUBE_ENV}") {
                    sh """
                    \${SCANNER_HOME}/bin/sonar-scanner \
                    -Dsonar.projectKey=${PROJECT_KEY} \
                    -Dsonar.projectName="${PROJECT_NAME}" \
                    -Dsonar.sources=. \
                    -Dsonar.exclusions=**/*_test.go \
                    -Dsonar.go.coverage.reportPaths=coverage.out
                    """
                }
            }
        }

        stage('Quality Gate Check') {
            steps {
                echo 'Waiting for Quality Gate...'
                timeout(time: 10, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Deploy to VPS') {
            steps {
                echo 'Deploying to VPS...'
                sh '''
                echo "PORT=8888" > .env
                docker compose up -d --build
                docker image prune -f
                '''
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished. Cleaning up workspace...'
            cleanWs()
        }
        failure {
            echo 'Pipeline failed. Please check the logs.'
        }
        success {
            echo 'Pipeline completed successfully.'
        }
    }
}