pipeline {
    agent any
    
    environment {
        PATH = "/usr/local/go/bin:$PATH"
        GOCACHE = "${WORKSPACE}/.cache/go-build"
        GOPATH = "${WORKSPACE}/go"
    }

    stages {
        stage('Build and Test (Parallel)') {
            parallel {
                stage('Unit Test') {
                    steps {
                        echo 'Running Unit Tests...'
                        sh 'go test -v ./...'
                    }
                }
                stage('Compile Build') {
                    steps {
                        echo 'Compiling the application...'
                        sh 'go build -o main .'
                    }
                }
            }
        }

        stage('SonarQube Analysis') {
            environment {
                SCANNER_HOME = tool 'sonar-scanner'
            }
            steps {
                echo 'Running SonarQube Analysis...'
                withSonarQubeEnv('SonarQube') {
                    sh """
                    \${SCANNER_HOME}/bin/sonar-scanner \
                    -Dsonar.projectKey=ncc-module-2 \
                    -Dsonar.projectName="NCC Module 2" \
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
                timeout(time: 5, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished. Cleaning up workspace...'
            cleanWs()
        }
    }
}