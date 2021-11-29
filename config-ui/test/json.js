var assert = require("assert");
const AVAILABLE_PLUGINS = require("../src/data/availablePlugins");
const TEST_DATA = require('./testData')
const { getCollectorJson, getCollectionJson } = require("../src/utils/triggersUtil");

describe("Json utils", () => {
  describe("getCollectionJson", function () {
    it("gets default JSON for plugins based on an array of names", function () {
      let expected = TEST_DATA.completeTriggersJson
      let actual = getCollectionJson(AVAILABLE_PLUGINS);
      assert.deepEqual(expected, actual);
    });
  });
  describe("getCollectorJson", function () {
    it("gets default JSON for a collector plugin based on the name", function () {
      const expected = TEST_DATA.gitlabTriggersJson
      let actual = getCollectorJson("gitlab");
      assert.deepEqual(expected, actual);
    });
  });
});
