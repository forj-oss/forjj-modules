pipeline {
    agent any

    stages {
        stage('prepare build environment') {
            steps {
                sh('''set +x ; source ./build-env.sh
                create-go-build-env.sh''')
            }
        }
        stage('Install dependencies') {
            steps {
                sh('''set +x ; source ./build-env.sh
                glide i''')
            }
        }
/*        stage('Building modules in parallel'){
            parallel {
                stage('Trace module') {
                    stages {*/
                        stage('Build trace'){
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go build forjj-modules/trace''')
                            }
                        }
                        stage('Tests trace') {
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go test forjj-modules/trace''')
                            }
                        }/*
                    }
                }
                stage('Trace cli') {
                    stages {*/
                        stage('Build cli') {
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go build forjj-modules/cli''')
                            }
                        }
                        stage('Tests cli') {
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go test forjj-modules/cli''')
                            }
                        }/*
                    }
                }
            }
        }*/
    }

    post {
        success {
            deleteDir()
        }
    }
}
