pipeline {
    agent any
    
    environment {
        DOCKER_BUILDKIT = '1'
        PATH = "/usr/local/go/bin:$PATH"
        GOCACHE = "${WORKSPACE}/.cache/go-build"
        GOMODCACHE = "${WORKSPACE}/go/pkg/mod"
        SONARQUBE_ENV = 'SonarQube'
        SCANNER_HOME = tool 'sonar-scanner'
        PROJECT_KEY   = 'ncc-module-2'
        PROJECT_NAME  = "NCC Module 2"
    }

    stages {
        stage('Test and Important Stuff') {
            parallel {
                stage('Unit Test') {
                    steps {
                        echo 'Running Unit Tests...'
                        sh 'go test -v ./... -coverprofile=coverage.out'
                    }
                }
                stage('Some important stage') {
                    steps {
                         echo 'Doing some very IMPORTANT stuffs...'
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

        stage('Build Docker Image') {
            steps {
                echo 'Building Docker Image...'
                sh 'docker build -t ncc-module-2:latest . || true'
            }
        }

        stage('Deploy to VPS') {
            steps {
                echo 'Deploying to VPS...'
                sh '''
                echo "PORT=8888" > .env
                docker compose pull || true
                docker compose up -d --build --force-recreate
                docker image prune -f
                '''
            }
        }
    }

    post {
        always {
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