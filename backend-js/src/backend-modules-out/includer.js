(() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __glob = (map) => (path) => {
    var fn = map[path];
    if (fn) return fn();
    throw new Error("Module not found in bundle: " + path);
  };
  var __esm = (fn, res) => function __init() {
    return fn && (res = (0, fn[__getOwnPropNames(fn)[0]])(fn = 0)), res;
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

  // ../../common/rpc.mts
  var methodHandlers;
  var init_rpc = __esm({
    "../../common/rpc.mts"() {
      methodHandlers = /* @__PURE__ */ new Map();
    }
  });

  // modules/files.mts
  var files_exports = {};
  var init_files = __esm({
    "modules/files.mts"() {
      init_rpc();
      methodHandlers["openFile"] = (name) => {
        return `Going to open file ${name}`;
      };
    }
  });

  // require("./modules/**/*.mts") in includer.js
  var globRequire_modules_mts = __glob({
    "./modules/files.mts": () => (init_files(), __toCommonJS(files_exports))
  });

  // includer.js
  var kind = "files";
  globRequire_modules_mts("./modules/" + kind + ".mts");
})();
