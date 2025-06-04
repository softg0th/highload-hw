@Library('COSM-Jenkins-libs') _

pipeline {
    agent none

    options {
        skipDefaultCheckout(true)
        timestamps()
    }

    environment {
        COMPOSE_PROJECT_NAME = "hl-${env.BUILD_TAG}"
        DOCKER_NETWORK = 'main_bridge'
        DOCKER_BUILDER_IMAGE = 'docker-builder'
        MINIO_HEALTHCHECK = 'http://localhost:9002/minio/health/ready'
    }

    stages {
        stage('Checkout & Preparation') {
            agent { label 'master' }
            steps {
                cleanWs()
                checkout scm
                sh '''
                    git log -n 1 | grep "commit " | sed 's/commit //g' > currenntVersion
                '''
                stash name: 'workspace', includes: '**'
            }
        }

        stage('Build Application') {
            agent {
                docker {
                    image 'maven:3.9.4-eclipse-temurin-17-alpine'
                    reuseNode true
                    args '-u root'
                }
            }
            steps {
                unstash 'workspace'
                sh '''
                    echo "[build] Put your real build commands here, e.g. mvn clean install"
                '''
            }
        }

        stage('Run Unit Tests') {
            agent {
                docker {
                    image 'golang:1.23'
                    reuseNode true
                    args '-u root'
                }
            }
            steps {
                unstash 'workspace'

                dir('filter') {
                    sh '''
                        echo "[go test] Running unit tests in filter service..."
                        go mod tidy
                        go test ./internal/core/... -v -coverprofile=coverage.out
                    '''
                    sh 'sleep 2'
                }

                dir('receiver') {
                    sh '''
                        echo "[go test] Running unit tests in receiver service..."
                        go mod tidy
                        go test ./internal/api/... -v -coverprofile=coverage.out
                    '''
                    sh 'sleep 2'
                }

                dir('storage') {
                    sh '''
                        echo "[go test] Running unit tests in storage service..."
                        go mod tidy
                        go test ./internal/api/... -v -coverprofile=coverage.out
                    '''
                    sh 'sleep 2'
                }
            }
            post {
                always {
                    junit allowEmptyResults: true, testResults: '**/{filter,receiver,storage}/**/TEST-*.xml'
                    archiveArtifacts artifacts: '**/coverage.out', allowEmptyArchive: true
                }
                failure {
                    echo '[go test] Unit tests failed'
                }
            }
        }



        stage('Start Docker Services') {
            agent {
                docker {
                    image "${DOCKER_BUILDER_IMAGE}"
                    reuseNode true
                    args "-u root --net=${DOCKER_NETWORK} -v /var/run/docker.sock:/var/run/docker.sock"
                }
            }
            steps {
                unstash 'workspace'

                sh '''
                    echo "[cleanup] Removing previous containers"
                    docker compose down -v --remove-orphans || true

                    echo "[cleanup] Removing leftover containers"
                    docker ps -aq --filter "name=highload" | xargs -r docker rm -f || true

                    echo "[docker] Start all containers except pytest/tsung"
                    docker compose up -d \
                        elasticsearch logstash kibana zookeeper kafka \
                        minio mongo prometheus grafana node_exporter \
                        storage receiver filter

                    echo "[docker] Waiting for containers to be healthy..."
                    sleep 10
                '''
            }
        }

        stage('Run Integration Tests') {
            agent {
                docker {
                    image "${DOCKER_BUILDER_IMAGE}"
                    reuseNode true
                    args "-u root --net=${DOCKER_NETWORK} -v /var/run/docker.sock:/var/run/docker.sock"
                }
            }
            steps {
                catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
                    sh '''
                        echo "[pytest] Deleting the old marker"
                        rm -f ./shared_tmp/pytest_done || true

                        echo "[pytest] Launching the test container"
                        docker compose up -d pytest

                        echo "[pytest] Waiting for the tests to be completed..."
                        timeout=180
                        elapsed=0
                        while [ ! -f ./shared_tmp/pytest_done ] && [ $elapsed -lt $timeout ]; do
                            echo "[pytest] still running..."
                            sleep 5
                            elapsed=$((elapsed + 5))
                        done

                        if [ ! -f ./shared_tmp/pytest_done ]; then
                            echo "[pytest] The waiting time has expired!"
                            exit 1
                        fi

                        echo "[pytest] Completed successfully"
                        docker logs ${COMPOSE_PROJECT_NAME}-pytest > pytest.log || echo "No logs pytest"
                    '''
                }
            }
            post {
                always {
                    archiveArtifacts artifacts: 'pytest.log', allowEmptyArchive: true
                }
                failure {
                    echo '[pytest] Integration tests failed.'
                }
            }
        }



//         stage('Run Load Tests') {
//             when {
//                 expression { fileExists('./shared_tmp/pytest_done') }
//             }
//             agent {
//                 docker {
//                     image "${DOCKER_BUILDER_IMAGE}"
//                     reuseNode true
//                     args "-u root --net=${DOCKER_NETWORK} -v /var/run/docker.sock:/var/run/docker.sock"
//                 }
//             }
//             steps {
//                 sh '''
//                     echo "[tsung] Starting load test container"
//                     docker compose up -d tsung
//
//                     echo "[tsung] Waiting for test to complete..."
//                     sleep 60 # You can improve this by checking logs or custom markers
//
//                     mkdir -p tsung_results
//                     cp -r ./tsung/log/* ./tsung_results/ || true
//                 '''
//             }
//             post {
//                 always {
//                     archiveArtifacts artifacts: 'tsung_results/**/*', allowEmptyArchive: true
//                 }
//                 failure {
//                     echo '[tsung] Load testing failed'
//                 }
//             }
//         }

        stage('Deploy Artifacts') {
            when {
                expression { fileExists('./shared_tmp/pytest_done') }
            }
            agent {
                docker {
                    image "${DOCKER_BUILDER_IMAGE}"
                    reuseNode true
                    args "-u root --net=${DOCKER_NETWORK} -v /var/run/docker.sock:/var/run/docker.sock"
                }
            }
            steps {
                sh '''
                    GIT_REVISION=$(cat currenntVersion)
                    echo "[deploy] Building final images"
                    docker compose build

                    echo "[deploy] Starting all services"
                    docker compose up -d
                '''
            }
        }
    }

    post {
        always {
            node('master') {
                script {
                    env.GIT_URL = env.GIT_URL_1
                    notifyRocketChat(
                        channelName: 'dummy',
                        minioCredentialsId: 'jenkins-minio-credentials',
                        minioHostUrl: 'https://minio.cloud.cosm-lab.science'
                    )

                    withCredentials([string(credentialsId: 'CloudRushTlg-token', variable: 'TLG_TOKEN')]) {
                        notifyTelegram(
                            minioHostUrl: 'https://minio.cloud.cosm-lab.science',
                            botIdAndToken: env.TLG_TOKEN,
                            chatId: '-1002474884172',
                            threadId: '2'
                        )
                    }
                }
            }
        }
    }
}

