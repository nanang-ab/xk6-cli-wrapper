import cli from 'k6/x/cli-wrapper'

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 1,
  // A string specifying the total itreartion of the test run.
  iterations: 1,

  // The following section contains configuration options for execution of this
  // test script in Grafana Cloud.
  //
  // See https://grafana.com/docs/grafana-cloud/k6/get-started/run-cloud-tests-from-the-cli/
  // to learn about authoring and running k6 test scripts in Grafana k6 Cloud.
  //
  // ext: {
  //   loadimpact: {
  //     // The ID of the project to which the test is assigned in the k6 Cloud UI.
  //     // By default tests are executed in default project.
  //     projectID: "",
  //     // The name of the test in the k6 Cloud UI.
  //     // Test runs with the same name will be grouped.
  //     name: "cli.js"
  //   }
  // },

  // Uncomment this section to enable the use of Browser API in your tests.
  //
  // See https://grafana.com/docs/k6/latest/using-k6-browser/running-browser-tests/ to learn more
  // about using Browser API in your test scripts.
  //
  // scenarios: {
  //   // The scenario name appears in the result summary, tags, and so on.
  //   // You can give the scenario any name, as long as each name in the script is unique.
  //   ui: {
  //     // Executor is a mandatory parameter for browser-based tests.
  //     // Shared iterations in this case tells k6 to reuse VUs to execute iterations.
  //     //
  //     // See https://grafana.com/docs/k6/latest/using-k6/scenarios/executors/ for other executor types.
  //     executor: 'shared-iterations',
  //     options: {
  //       browser: {
  //         // This is a mandatory parameter that instructs k6 to launch and
  //         // connect to a chromium-based browser, and use it to run UI-based
  //         // tests.
  //         type: 'chromium',
  //       },
  //     },
  //   },
  // }
};

export function setup() {
  const downloadUrl = "https://cdn.development.armada.accelbyte.io/linux_amd64/ams";
  // const expectedHash = "827b1dd34a403c30699b176f31e231ca10aa55dc6665ed6bec62d6c7b03f7bab";
  const cliDirectory = cli.getWritableDirectory(".run/files/ams/cli");

  const cliPath = cli.downloadCLIFile(downloadUrl, cliDirectory);
  console.log("downloaded cliPath: " + cliPath);

  // TODO: uncomment this after we have official hash sum
  // let isValid = cli.validateCLIFileHash(cliPath, expectedHash);
  // if (isValid !== null && !isValid) {
  //   throw new Error("CLI app validation failed: " + isValid);
  // }

  const absCliPath = cli.getAbsolutePath(cliPath);
  console.log(`CLI app downloaded at : [${absCliPath}] is valid and ready for use.`);

  return { cliPath: cliPath };
}

export function teardown(data) {
    cli.cleanupCLI(data.cliPath)
}

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function(data) {
  let args = [
    "upload",
     "-p" , ".run/files/ams/mockds",
     "-H", "development.accelbyte.io",
     "-c", "<client-id>", // replace it with client id
     "-s" ,"<client-secret>", // replace it with client secret
     "-n", "mockds-image",
     "-e", "mockds",
  ];

  try {
    let res = cli.executeCommand(data.cliPath, ...args);
    console.info(`Command execution success: ${res}`)
  } catch (error) {
    console.error(`Command execution error: ${error}`)
  }

}
