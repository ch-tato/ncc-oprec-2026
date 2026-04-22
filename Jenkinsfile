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
                        sh 'go test -v ./... -coverprofile=coverage.out || true'
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
                sh '''
                if docker image inspect ncc-module-2:latest > /dev/null 2>&1; then
                    docker build \
                        --build-arg BUILDKIT_INLINE_CACHE=1 \
                        --cache-from=ncc-module-2:latest \
                        -t ncc-module-2:latest .
                else
                    docker build -t ncc-module-2:latest .
                fi
                '''
            }
        }

        stage('Deploy to VPS') {
            steps {
                echo 'Deploying to VPS...'
                sh '''
                echo "PORT=8888" > .env
                docker compose up -d
                '''
            }
        }
    }

    post {
        // always {
        //     echo 'Pipeline finished. Cleaning up workspace...'
        //     cleanWs patterns: [[pattern: '.cache/**', type: 'EXCLUDE']]
        // }
        failure {
            echo 'Pipeline failed. Please check the logs.'
        }
        success {
            echo 'Pipeline completed successfully.'
        }
    }
}