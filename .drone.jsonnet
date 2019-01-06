local DRONE_CLI_VERSION = '1.0.4';
local DRONE_TOKEN = '55f24eb3d61ef6ac5e83d550178638dc';

local terraform_versions = [
  '0.10.6',
  '0.10.7',
  '0.10.8',
  '0.11.9',
  '0.11.10',
  '0.11.11',
];

local buildENV() = {
  CGO_ENABLE: '0',
  GO111MODULE: 'on',
  DRONE_TOKEN: DRONE_TOKEN,
  DRONE_SERVER: 'http://drone',
};


local stepTest() = {
  name: 'Unit Test',
  image: 'golang:1.11-alpine',
  environment: buildENV(),
  commands: [
    'apk add build-base',
    'go test -mod=vendor ./drone',
  ],
};

local pipelineTest() = {
  kind: 'pipeline',
  name: 'Unit Test',
  steps: [
    stepTest(),
  ],
};


local stepBuild() = {
  name: 'Build Drone Provider',
  image: 'golang:1.11-alpine',
  environment: buildENV(),
  commands: [
    'apk add build-base',
    'go build -mod=vendor',
    'cp terraform-provider-drone ci/terraform.d/plugins/linux_amd64/terraform-provider-drone_v0.0.0',
  ],
};

// setup gitea and drone sevices
local stepSetupServices() = {
  name: 'Setup Services',
  image: 'alpine',
  commands: [
    'apk add --update curl',
    // wait for gitea and wait for the restart
    'until nc -z gitea 3000; do sleep 5; done',
    'sleep 5',
    'until nc -z gitea 3000; do sleep 5; done',
    // create a test repo
    'until curl -f -u test:test -X POST -H "Content-Type: application/json" -d \'{"auto_init":true,"name":"test","readme":"Default"}\' http://gitea:3000/api/v1/user/repos; do sleep 5; done',
    // login to drone, drone sets up the connection gitea
    'sleep 5',
    "until curl -f -X POST -d 'username=test&password=test' http://drone/login; do sleep 5; done",
  ],
};

local stepIntegrationRepoCreate(tf_version) = {
  name: 'Test Repo Create',
  image: 'hashicorp/terraform:' + tf_version,
  environment: buildENV(),
  commands: [
    'apk add --update curl jq',
    'cd ci',
    'cp create/repo.tf ./',
    'terraform init',
    'terraform apply -auto-approve',
    // test if the repo was activated
    'curl -sSf http://drone/api/repos/test/test | jq .',
    "curl -sSf http://drone/api/repos/test/test | jq '.active|contains(true)'",
  ],
};

local stepIntegrationRepoUpdate(tf_version) = {
  name: 'Test Repo Update',
  image: 'hashicorp/terraform:' + tf_version,
  environment: buildENV(),
  commands: [
    'apk add --update curl jq',
    'cd ci',
    'cp update/repo.tf ./',
    'terraform apply -auto-approve',
    // test if the repo was activated
    'curl -sSf http://drone/api/repos/test/test | jq .',
    "curl -sSf http://drone/api/repos/test/test | jq '.trusted|contains(true)'",
  ],
};

local serviceGitea() = {
  name: 'gitea',
  image: 'gitea/gitea:latest',
  commands: [
    's6-svscan /etc/s6 &',
    'apk add sed',
    // wait for gitea to generate the inital config
    'until nc -z localhost 3000; do sleep 5; done',
    'until stat /data/gitea/conf/app.ini; do sleep 5; done',
    // set gitea as installed
    @"sed -i -e 's/INSTALL_LOCK\\s*=\\s*false/INSTALL_LOCK = true/' /data/gitea/conf/app.ini",
    // disable sql logs
    @"sed -i -e 's/\\[database\\]/[database]\\nLOG_SQL = false/' /data/gitea/conf/app.ini",
    // restart gitea
    'until killall gitea; do sleep 5; done',
    'until nc -z localhost 3000; do sleep 5; done',
    // create a test user
    'until gitea admin create-user --name test --password test --email test@test.local --admin --config /data/gitea/conf/app.ini; do sleep 5; done',
    'wait',
  ],
};

local serviceDrone() = {
  name: 'drone',
  image: 'drone/drone:1.0.0-rc.3',
  environment: {
    DRONE_GITEA_SERVER: 'http://gitea:3000',
    DRONE_GIT_USER: 'test',
    DRONE_GIT_PASSWORD: 'test',
    DRONE_RUNNER_CAPACITY: '1',
    DRONE_SERVER_HOST: 'localhost',
    DRONE_SERVER_PROTO: 'http',
    DRONE_TLS_AUTOCERT: 'false',
    DRONE_USER_CREATE: 'username:test,machine:false,admin:true,token:' + DRONE_TOKEN,
  },
};

local pipelineIntegration(tf_version) = {
  kind: 'pipeline',
  name: 'Integration Test ' + tf_version,
  steps: [
    stepSetupServices(),
    stepBuild(),
    stepIntegrationRepoCreate(tf_version),
    stepIntegrationRepoUpdate(tf_version),
  ],
  services: [
    serviceGitea(),
    serviceDrone(),
  ],
  depends_on: [
    'Unit Test',
  ],

};

// local stepBuild(os, arch) = {
// name: 'build',
// image: 'golang:alpine-1.11',
// environment: {
// CGO_ENABLE: 0,
// GOOS
// },
// commands: [
// std.format('GOOS=%(os)s GOARCH=%(arch)s go build -mod=vendor dist/%(os)s_%(arch)s/terraform-provider-drone_${DRONE_TAG}',
// {os: os, arch: arch}),
// std.format('tar -cvzf dist/terraform-provider-drone_%(os)s_%(arch)s.tar.gz -C dist/%(os)s_%(arch)s terraform-provider')
// ]
// };

[
  pipelineTest(),
] + [pipelineIntegration(tfv) for tfv in terraform_versions]
