pipeline {
    agent any

    stages {
        stage('prepare build environment') {
            steps {
                sh('''set +x ; source ./build-env.sh
                create-go-build-env.sh''')
            }
        }
        stage('Check self GO module reference') {
            steps {
                sh('''set +x ; source ./build-env.sh
                bin/go-check_files.sh
                ''')
            }
        }
        stage('Install dependencies') {
            steps {
                sh('''set +x ; source ./build-env.sh
                glide i''')
            }
        }
        stage('Building modules in parallel'){
            parallel {
                stage('Trace module') {
                    stages {
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
                        }
                    }
                }
                stage('cli module') {
                    stages {
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
                        }
                    }
                }
                stage('cli/kingpinCli module') {
                    stages {
                        stage('Build cli/kingpinCli') {
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go build forjj-modules/cli/kingpinCli''')
                            }
                        }
                        stage('Tests cli') {
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go test forjj-modules/cli/kingpinCli''')
                            }
                        }
                    }
                }
                stage('cli/kingpinMock module') {
                    stages {
                        stage('Build cli/kingpinMock') {
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go build forjj-modules/cli/kingpinMock''')
                            }
                        }
                        stage('Tests cli') {
                            steps {
                                sh('''set +x ; source ./build-env.sh
                                go test forjj-modules/cli/kingpinMock''')
                            }
                        }
                    }
                }
            }
        }
    }

    post {
        success {
            deleteDir()
        }
    }
}
