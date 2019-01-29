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
                sh('''set +x ; source ./build-env.sh ; set -x
                if [[ "$(find . -wholename ./vendor -prune -o -wholename ./.glide -prune -o -name "*.go" -exec grep \\"$BE_PROJECT/ {} \\; | wc -l)" -ne 0 ]]
                then
                    echo "A GO module requires self reference to be exported. (relative path is not accepted)"
                    exit 1
                fi
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
                stage('Trace cli') {
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
            }
        }
    }

    post {
        success {
            deleteDir()
        }
    }
}
