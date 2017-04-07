/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId])
/******/ 			return installedModules[moduleId].exports;
/******/
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// identity function for calling harmony imports with the correct context
/******/ 	__webpack_require__.i = function(value) { return value; };
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, {
/******/ 				configurable: false,
/******/ 				enumerable: true,
/******/ 				get: getter
/******/ 			});
/******/ 		}
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = 22);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var log_1 = __webpack_require__(1);
function asArray(collection) {
    return Array.prototype.slice.apply(collection);
}
exports.asArray = asArray;
function select(w, r) {
    var sel = w.getSelection();
    sel.removeAllRanges();
    sel.addRange(r);
}
exports.select = select;
function getSelectedRange(w) {
    try {
        return w.getSelection().getRangeAt(0);
    }
    catch (ignore) {
        return null;
    }
}
exports.getSelectedRange = getSelectedRange;
function clearSelection(w) {
    w.getSelection().removeAllRanges();
}
exports.clearSelection = clearSelection;
/**
 * Find the first element containing n that can be used as a position anchor.
 */
function findPositionAnchor(w, n) {
    var e = null;
    if (n.nodeType != Node.ELEMENT_NODE) {
        e = n.parentElement;
    }
    else {
        e = n;
    }
    if (e == null) {
        throw 'no element parent found';
    }
    do {
        var stl = w.getComputedStyle(e);
        if (stl.position == "relative" || stl.position == "fixed") {
            return e;
        }
        else {
            e = e.parentElement;
        }
    } while (e);
    var b = w.document.body;
    log_1.warn("no position anchor found, using body (" + w.getComputedStyle(b).position + ").");
    return b;
}
exports.findPositionAnchor = findPositionAnchor;
function computeUniquePath(n) {
    var path = '';
    if (n.nodeType == Node.ELEMENT_NODE) {
        var e = n;
        if (e.id) {
            return "#" + e.id;
        }
    }
    if (n.nodeName == 'BODY') {
        return '/';
    }
    var siblings = asArray(n.parentNode.childNodes);
    var idx = siblings.indexOf(n);
    return computeUniquePath(n.parentNode) + ("/" + idx);
}
exports.computeUniquePath = computeUniquePath;
function getNodeByPath(doc, uPath) {
    var n = doc.body;
    var parts = uPath.split('/');
    for (var _i = 0, parts_1 = parts; _i < parts_1.length; _i++) {
        var p = parts_1[_i];
        if (!p.trim()) {
            continue;
        }
        if (p.indexOf('#') == 0) {
            n = doc.getElementById(p.substring(1));
            continue;
        }
        n = n.childNodes[parseInt(p)];
    }
    return n;
}
exports.getNodeByPath = getNodeByPath;
var BoxSizing;
(function (BoxSizing) {
    BoxSizing[BoxSizing["ContentBox"] = 0] = "ContentBox";
    BoxSizing[BoxSizing["BorderBox"] = 1] = "BorderBox";
})(BoxSizing = exports.BoxSizing || (exports.BoxSizing = {}));
function getBoxSizing(w, e) {
    var bs = w.getComputedStyle(e).boxSizing;
    if (bs == 'content-box') {
        return BoxSizing.ContentBox;
    }
    else if (bs = 'border-box') {
        return BoxSizing.BorderBox;
    }
    else {
        throw 'unknown box sizing option.';
    }
}
exports.getBoxSizing = getBoxSizing;
function getStyleNumber(stl, name) {
    var strVal = stl[name];
    var match = /\d+/.exec(strVal);
    if (match) {
        return Number(match[0]);
    }
    throw Error("cannot parse number from style " + name + ": " + strVal);
}
exports.getStyleNumber = getStyleNumber;


/***/ }),
/* 1 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var LogLevel;
(function (LogLevel) {
    LogLevel[LogLevel["DEBUG"] = 0] = "DEBUG";
    LogLevel[LogLevel["INFO"] = 1] = "INFO";
    LogLevel[LogLevel["WARN"] = 2] = "WARN";
    LogLevel[LogLevel["ERROR"] = 3] = "ERROR";
})(LogLevel = exports.LogLevel || (exports.LogLevel = {}));
var __level = LogLevel.DEBUG;
function setLogLevel(level) {
    __level = level;
}
exports.setLogLevel = setLogLevel;
function levelToString(level) {
    switch (level) {
        case LogLevel.DEBUG:
            return 'debug';
        case LogLevel.INFO:
            return 'info';
        case LogLevel.WARN:
            return 'warn';
        case LogLevel.ERROR:
            return 'error';
        default:
            return 'debug';
    }
}
function log(level, messages) {
    if (level >= __level) {
        if (console) {
            var stub = levelToString(level);
            if (messages.length <= 0) {
                messages = ['\r\n'];
            }
            var msg = messages.length > 1 ? messages : messages[0];
            try {
                (console[stub])(msg);
            }
            catch (ignore) { }
        }
    }
}
function debug() {
    var messages = [];
    for (var _i = 0; _i < arguments.length; _i++) {
        messages[_i] = arguments[_i];
    }
    log(LogLevel.DEBUG, messages);
}
exports.debug = debug;
function info() {
    var messages = [];
    for (var _i = 0; _i < arguments.length; _i++) {
        messages[_i] = arguments[_i];
    }
    log(LogLevel.INFO, messages);
}
exports.info = info;
function warn() {
    var messages = [];
    for (var _i = 0; _i < arguments.length; _i++) {
        messages[_i] = arguments[_i];
    }
    log(LogLevel.WARN, messages);
}
exports.warn = warn;
function error() {
    var messages = [];
    for (var _i = 0; _i < arguments.length; _i++) {
        messages[_i] = arguments[_i];
    }
    log(LogLevel.ERROR, messages);
}
exports.error = error;


/***/ }),
/* 2 */
/***/ (function(module, exports) {

module.exports = React;

/***/ }),
/* 3 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
function MeasureSpanClass() {
    return 'hlc_measure_span';
}
exports.MeasureSpanClass = MeasureSpanClass;
function DefaultColor() {
    return "green";
}
exports.DefaultColor = DefaultColor;
function DefaultFocusColor() {
    return "green";
}
exports.DefaultFocusColor = DefaultFocusColor;
function DefaultOpacity() {
    return 0.25;
}
exports.DefaultOpacity = DefaultOpacity;
function DefaultZIndex() {
    return 99999999;
}
exports.DefaultZIndex = DefaultZIndex;
function SameRowHFactor() {
    return 8;
}
exports.SameRowHFactor = SameRowHFactor;
function SameRowVFactor() {
    return 0.25;
}
exports.SameRowVFactor = SameRowVFactor;


/***/ }),
/* 4 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var range_meta_1 = __webpack_require__(19);
var dom_helper_1 = __webpack_require__(0);
/**
 * DOM Range object will change if DOM tree changed.
 * RangeCache will not.
 * You should use RangeCache.make(doc, range) method to creat an
 * intance.
 */
var RangeCache = (function () {
    function RangeCache(doc, cac, start, end, startOffset, endOffset, meta) {
        _a = [doc, cac, start, end, startOffset, endOffset], this.m_document = _a[0], this.m_cac = _a[1], this.m_start = _a[2], this.m_end = _a[3], this.m_startOffset = _a[4], this.m_endOffset = _a[5];
        this.setMeta(meta);
        var _a;
    }
    /**
     * A range expires when the dom sub tree is destroyed.
     */
    RangeCache.prototype.isExpired = function () {
        for (var _i = 0, _a = [this.m_cac, this.m_start, this.m_end]; _i < _a.length; _i++) {
            var n = _a[_i];
            if (!n || !n.parentNode) {
                return true;
            }
        }
        return false;
    };
    RangeCache.prototype.toRange = function (doc) {
        var r = doc.createRange();
        r.setStart(this.m_start, this.m_startOffset);
        r.setEnd(this.m_end, this.m_endOffset);
        return r;
    };
    Object.defineProperty(RangeCache.prototype, "document", {
        get: function () {
            return this.m_document;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeCache.prototype, "commonAncestorContainer", {
        get: function () {
            return this.m_cac;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeCache.prototype, "startContainer", {
        get: function () {
            return this.m_start;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeCache.prototype, "endContainer", {
        get: function () {
            return this.m_end;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeCache.prototype, "startOffset", {
        get: function () {
            return this.m_startOffset;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeCache.prototype, "endOffset", {
        get: function () {
            return this.m_endOffset;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeCache.prototype, "meta", {
        get: function () {
            return this.m_meta;
        },
        enumerable: true,
        configurable: true
    });
    RangeCache.prototype.setMeta = function (meta) {
        if (!meta) {
            var anchor = dom_helper_1.findPositionAnchor(this.m_document.defaultView, this.m_cac);
            this.m_meta = new range_meta_1.RangeMeta([
                dom_helper_1.computeUniquePath(anchor),
                dom_helper_1.computeUniquePath(this.m_start),
                dom_helper_1.computeUniquePath(this.m_end),
                this.toRange(this.m_document).toString()
            ], [this.m_startOffset, this.m_endOffset]);
        }
        else {
            this.m_meta = meta;
        }
    };
    return RangeCache;
}());
RangeCache.make = function (doc, r, meta) {
    return new RangeCache(doc, r.commonAncestorContainer, r.startContainer, r.endContainer, r.startOffset, r.endOffset, meta);
};
exports.RangeCache = RangeCache;


/***/ }),
/* 5 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var node_measure_1 = __webpack_require__(6);
/**
 * A block of a text highlighted in one selection
 */
var Block = (function () {
    function Block(id, rangeCache, dimensions) {
        this.m_id = id;
        this.m_rangeCache = rangeCache;
        this.m_dimensions = dimensions;
    }
    Block.prototype.setId = function (id) {
        this.m_id = id;
    };
    Object.defineProperty(Block.prototype, "id", {
        get: function () {
            return this.m_id;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Block.prototype, "rangeCache", {
        get: function () {
            return this.m_rangeCache;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Block.prototype, "rangeMeta", {
        get: function () {
            return this.m_rangeCache.meta;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Block.prototype, "dimensions", {
        get: function () {
            return this.m_dimensions;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Block.prototype, "text", {
        get: function () {
            return this.rangeMeta.text;
        },
        enumerable: true,
        configurable: true
    });
    /**
     * recalculateDimension
     * Recalculate dimensions of the range again.
     */
    Block.prototype.recalculateDimension = function () {
        this.m_dimensions =
            node_measure_1.computeDimentions(this.m_rangeCache.document, this.m_rangeCache);
    };
    return Block;
}());
exports.Block = Block;


/***/ }),
/* 6 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var range_cache_1 = __webpack_require__(4);
var dom_helper_1 = __webpack_require__(0);
var array_helper_1 = __webpack_require__(20);
var log_1 = __webpack_require__(1);
var node_measure_impl_1 = __webpack_require__(17);
var range_helper_1 = __webpack_require__(18);
function computeDimentions(doc, range) {
    var dims = [];
    var rc = null;
    if (range_helper_1.isRange(range)) {
        rc = range_cache_1.RangeCache.make(doc, range);
    }
    else {
        rc = range;
    }
    rc = range_helper_1.correctRange(rc);
    range_helper_1.iterateRangeNodes(rc, buildNodeToDimVisitor(document, dims));
    dims = node_measure_impl_1.mergeDimensions(dims);
    log_1.debug('computed dimensions: ');
    dims.forEach(function (v) {
        log_1.debug(v.toString());
    });
    log_1.debug('- - - - - - - - - -');
    return dims;
}
exports.computeDimentions = computeDimentions;
function buildNodeToDimVisitor(doc, result) {
    return function (n, ctx, s, e) {
        var dims = nodeToDimensions(doc, n, ctx, s, e);
        array_helper_1.extend(result, dims);
    };
}
function nodeToDimensions(doc, n, ctx, start, end) {
    if (n.nodeType == Node.TEXT_NODE) {
        return textNodeToDim(doc, n, ctx, start, end);
    }
    else {
        return [];
    }
}
function textNodeToDim(doc, n, ctx, start, end) {
    var span = node_measure_impl_1.substitudeWithMeasureSpan(doc, n, ctx, start, end);
    var dims = node_measure_impl_1.measureSpanToDim(span, doc, function (doc, el) {
        var anchor = dom_helper_1.getNodeByPath(doc, ctx.rangeCache.meta.anchorUPath);
        return node_measure_impl_1.computeLayout(doc, anchor, el);
    });
    node_measure_impl_1.restoreBeforeMeasureStatus(ctx.parent, ctx.index, ctx.siblings);
    return dims;
}


/***/ }),
/* 7 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __extends = (this && this.__extends) || (function () {
    var extendStatics = Object.setPrototypeOf ||
        ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
        function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
var log_1 = __webpack_require__(1);
var Item_1 = __webpack_require__(13);
var React = __webpack_require__(2);
var Footnote = (function (_super) {
    __extends(Footnote, _super);
    function Footnote() {
        var _this = _super.call(this) || this;
        _this.state = { highlights: [] };
        return _this;
    }
    Footnote.prototype.componentDidMount = function () {
        var _this = this;
        if (this.props.source) {
            this.addBatch(this.props.source.get());
            if (this.props.source.poll) {
                setInterval(function () {
                    _this.addBatch(_this.props.source.get());
                }, this.props.source.interval);
            }
        }
        this.props.onMount(this);
    };
    Footnote.prototype.addHighlight = function (h) {
        this.addBatch([h]);
    };
    Footnote.prototype.addBatch = function (hs) {
        var hlts = this.state.highlights.slice();
        for (var _i = 0, hs_1 = hs; _i < hs_1.length; _i++) {
            var h = hs_1[_i];
            hlts.push(h);
        }
        this.setState({ highlights: hlts });
    };
    Footnote.prototype.removeHighlight = function (id) {
        var hlts = this.state.highlights.slice();
        var found = false;
        for (var idx = 0; idx < hlts.length; idx++) {
            var h = hlts[idx];
            if (h.id == id) {
                hlts.splice(idx, 1);
                found = true;
                this.setState({ highlights: hlts });
                break;
            }
        }
        return found;
    };
    Footnote.prototype.removeAll = function () {
        this.setState({ highlights: [] });
    };
    Footnote.prototype.getAll = function () {
        return this.state.highlights.slice();
    };
    Footnote.prototype.onHighlightClick = function (i, e) {
        log_1.debug(i + " was clicked.");
    };
    Footnote.prototype.render = function () {
        var _this = this;
        return React.createElement("div", { id: this.props.id }, this.state.highlights.map(function (blk, idx) {
            return React.createElement(Item_1.Item, { key: blk.id, block: blk, onClick: function (e) { return _this.onHighlightClick(idx, e); } });
        }));
    };
    return Footnote;
}(React.Component));
exports.Footnote = Footnote;


/***/ }),
/* 8 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var range_cache_1 = __webpack_require__(4);
var dom_helper_1 = __webpack_require__(0);
var id_helper_1 = __webpack_require__(21);
var log_1 = __webpack_require__(1);
var block_1 = __webpack_require__(5);
var node_measure_1 = __webpack_require__(6);
/**
 * Detect if the given window has a expanded selection range.
 * Return a Block object containing the dimension information
 * of the range or null.
 *
 * @param win window object to extract selection
 * @param doc document object of the window
 * @param id [optional] if present, use it as the
 * Block object's id; otherwise generate a random uuid.
 */
function extractSelectedBlock(win, doc, id) {
    var range = dom_helper_1.getSelectedRange(win);
    if (null == range || range.collapsed) {
        log_1.debug('no selected range detected.');
        return null;
    }
    var rc = range_cache_1.RangeCache.make(doc, range);
    dom_helper_1.clearSelection(win);
    return generateBlock(doc, rc, id);
}
exports.extractSelectedBlock = extractSelectedBlock;
/**
 * Rebuild a block using given metadata.
 *
 * @param win window object to extract selection
 * @param doc document object of the window
 * @param meta metadata used to restore the block. See @class{RangeMeta}
 * @param id [optional] if present, use it as the
 * Block object's id; otherwise generate a random uuid.
 */
function restoreBlock(win, doc, meta, id) {
    var rc = restoreRangeCache(doc, meta);
    return generateBlock(doc, rc, id);
}
exports.restoreBlock = restoreBlock;
function generateBlock(doc, rangeCache, id) {
    var dims = node_measure_1.computeDimentions(doc, rangeCache);
    if (!id) {
        id = id_helper_1.generateRandomUUID();
    }
    return new block_1.Block(id, rangeCache, dims);
}
function restoreRangeCache(doc, meta) {
    var rangeAnchors = [];
    for (var _i = 0, _a = [meta.startNodeUPath, meta.endNodeUPath]; _i < _a.length; _i++) {
        var uPath = _a[_i];
        try {
            var n = dom_helper_1.getNodeByPath(doc, uPath);
            rangeAnchors.push(n);
        }
        catch (e) {
            log_1.warn(e);
            return null;
        }
    }
    var r = doc.createRange();
    r.setStart(rangeAnchors[0], meta.startCharIndex);
    r.setEnd(rangeAnchors[1], meta.endCharIndex);
    var rc = range_cache_1.RangeCache.make(doc, r, meta);
    return rc;
}


/***/ }),
/* 9 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __extends = (this && this.__extends) || (function () {
    var extendStatics = Object.setPrototypeOf ||
        ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
        function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
var block_1 = __webpack_require__(5);
var DecoratedBlock = (function (_super) {
    __extends(DecoratedBlock, _super);
    function DecoratedBlock(id, rc, m, s) {
        var _this = _super.call(this, id, rc, m) || this;
        _this.m_styles = s;
        return _this;
    }
    Object.defineProperty(DecoratedBlock.prototype, "styles", {
        get: function () {
            return this.m_styles;
        },
        enumerable: true,
        configurable: true
    });
    DecoratedBlock.prototype.decorate = function (decorator) {
        for (var i = 0; i < this.dimensions.length; i++) {
            var dim = this.dimensions[i];
            var stl = this.m_styles[i];
            decorator(dim, stl);
        }
    };
    return DecoratedBlock;
}(block_1.Block));
exports.DecoratedBlock = DecoratedBlock;
var DecoratedBlockFactory = (function () {
    function DecoratedBlockFactory(block) {
        this.m_block = block;
        this.m_decorators = [];
    }
    DecoratedBlockFactory.prototype.addDecorator = function (d) {
        this.m_decorators.push(d);
        return this;
    };
    DecoratedBlockFactory.prototype.make = function () {
        var styles = [];
        for (var _i = 0, _a = this.m_block.dimensions; _i < _a.length; _i++) {
            var dim = _a[_i];
            var s = {};
            for (var _b = 0, _c = this.m_decorators; _b < _c.length; _b++) {
                var decorate = _c[_b];
                decorate(dim, s);
            }
            styles.push(s);
        }
        return new DecoratedBlock(this.m_block.id, this.m_block.rangeCache, this.m_block.dimensions, styles);
    };
    return DecoratedBlockFactory;
}());
exports.DecoratedBlockFactory = DecoratedBlockFactory;


/***/ }),
/* 10 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var constants_1 = __webpack_require__(3);
function BasicDecorator(dim, styles) {
    styles.position = dim.Fixed ? "fixed" : "absolute";
    styles.top = dim.Top;
    styles.left = dim.Left;
    styles.width = dim.Width;
    styles.height = dim.Height;
    styles.backgroundColor = constants_1.DefaultColor();
    styles.opacity = constants_1.DefaultOpacity();
    styles.zIndex = constants_1.DefaultZIndex();
    styles.pointerEvents = "none";
}
exports.BasicDecorator = BasicDecorator;


/***/ }),
/* 11 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/**
 * Unsafely generating markdown formatted text
 * https://guides.github.com/features/mastering-markdown/#syntax
 */

Object.defineProperty(exports, "__esModule", { value: true });
function h(level, text) {
    return repeat('#', level) + " " + text;
}
exports.h = h;
function em(text) {
    return "_" + text + "_";
}
exports.em = em;
function link(text, link) {
    return "[" + text + "](" + link + ")";
}
exports.link = link;
function ul(items, level) {
    if (level === void 0) { level = 0; }
    var md = "";
    for (var _i = 0, items_1 = items; _i < items_1.length; _i++) {
        var item = items_1[_i];
        md += repeat("  ", level) + "* " + item + "\r\n";
    }
    return md;
}
exports.ul = ul;
/**
 * Reference:
 * https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/repeat#Polyfill
 * @param text text to repeat
 * @param count time to repeat
 */
function repeat(text, count) {
    if (text == null) {
        throw new Error("cannot repeat null");
    }
    var str = "" + text;
    if (count < 0) {
        throw new Error("Cannot repeat negative times");
    }
    if (count == Infinity) {
        throw new Error("Cannot repeat infinite times");
    }
    count = Math.floor(count);
    if (str.length == 0 || count == 0) {
        return "";
    }
    if (str.length * count >= 1 << 28) {
        throw new RangeError("repeat count must not overflow maximum string size");
    }
    var rpt = '';
    for (;;) {
        if ((count & 1) == 1) {
            rpt += str;
        }
        count >>>= 1;
        if (count == 0) {
            break;
        }
        str += str;
    }
    return rpt;
}


/***/ }),
/* 12 */
/***/ (function(module, exports) {

module.exports = ReactDOM;

/***/ }),
/* 13 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var Row_1 = __webpack_require__(14);
var React = __webpack_require__(2);
exports.Item = function (props) {
    return React.createElement("div", { onClick: function (e) { return props.onClick; } }, props.block.styles.map(function (stl, i) { return React.createElement(Row_1.Row, { key: i, style: stl }); }));
};


/***/ }),
/* 14 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var React = __webpack_require__(2);
/**
 * Row represents one highlight unit
 */
exports.Row = function (props) {
    return React.createElement("div", { style: props.style, onClick: function () { return props.onClick; }, onMouseOver: function () { return props.onMouseOver; }, onMouseOut: function () { return props.onMouseOut; } });
};


/***/ }),
/* 15 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
/**
 * Measurements of a rectangular block.
 */
var Dimension = (function () {
    function Dimension(_a) {
        var x = _a[0], y = _a[1], w = _a[2], h = _a[3];
        _b = [x, y, w, h], this.m_x = _b[0], this.m_y = _b[1], this.m_width = _b[2], this.m_height = _b[3];
        this.m_fixed = false;
        var _b;
    }
    Object.defineProperty(Dimension.prototype, "Left", {
        get: function () {
            return this.m_x;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Dimension.prototype, "Top", {
        get: function () {
            return this.m_y;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Dimension.prototype, "Width", {
        get: function () {
            return this.m_width;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Dimension.prototype, "Height", {
        get: function () {
            return this.m_height;
        },
        enumerable: true,
        configurable: true
    });
    Dimension.prototype.setFixed = function () {
        this.m_fixed = true;
    };
    Object.defineProperty(Dimension.prototype, "Fixed", {
        get: function () {
            return this.m_fixed;
        },
        enumerable: true,
        configurable: true
    });
    Dimension.prototype.toString = function () {
        var d = this;
        return "{(" + d.Left + ", " + d.Top + "), " + d.Width + " x " + d.Height + "}";
    };
    return Dimension;
}());
exports.Dimension = Dimension;


/***/ }),
/* 16 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var dom_helper_1 = __webpack_require__(0);
/**
 * Store the context information of a node. This is useful if you
 * are about to temporarily take the node off the dom tree / fragment.
 */
var NodeContext = (function () {
    function NodeContext(n, rc) {
        this.m_parent = n.parentNode;
        this.m_nextSibling = n.nextSibling;
        this.m_siblings = dom_helper_1.asArray(n.parentNode.childNodes);
        this.m_rangeCache = rc;
        this.m_index = this.m_siblings.indexOf(n);
    }
    Object.defineProperty(NodeContext.prototype, "parent", {
        get: function () {
            return this.m_parent;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(NodeContext.prototype, "nextSibling", {
        get: function () {
            return this.m_nextSibling;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(NodeContext.prototype, "siblings", {
        get: function () {
            return this.m_siblings;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(NodeContext.prototype, "index", {
        get: function () {
            return this.m_index;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(NodeContext.prototype, "rangeCache", {
        get: function () {
            return this.m_rangeCache;
        },
        enumerable: true,
        configurable: true
    });
    return NodeContext;
}());
exports.NodeContext = NodeContext;


/***/ }),
/* 17 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var dimension_1 = __webpack_require__(15);
var constants = __webpack_require__(3);
var constants_1 = __webpack_require__(3);
var dom_helper = __webpack_require__(0);
var dom_helper_1 = __webpack_require__(0);
function measureSpanToDim(s, doc, calc) {
    var cspans = dom_helper.asArray(s.children);
    var dims = [];
    for (var _i = 0, cspans_1 = cspans; _i < cspans_1.length; _i++) {
        var span = cspans_1[_i];
        var dim = calc(doc, span);
        dims.push(dim);
    }
    return dims;
}
exports.measureSpanToDim = measureSpanToDim;
function mergeDimensions(dims) {
    var merged = [];
    var current = null;
    var pre = null;
    var maxIdx = dims.length - 1;
    for (var i in dims) {
        var d = dims[i];
        if (current == null) {
            current = d;
            pre = d;
            continue;
        }
        var charW = pre.Width;
        var charH = pre.Height;
        if (isInsameRow(d, current, charW, charH)) {
            current = doMerge(d, current);
        }
        else {
            merged.push(current);
            current = d;
        }
        if (Number(i) == maxIdx) {
            merged.push(current);
        }
        else {
            pre = d;
        }
    }
    return merged;
}
exports.mergeDimensions = mergeDimensions;
function isInsameRow(a, b, charW, charH) {
    if (a.Left > b.Left) {
        _a = [b, a], a = _a[0], b = _a[1];
    }
    var maxHGap = charW * constants_1.SameRowHFactor();
    var maxVGap = charH * constants_1.SameRowVFactor();
    var h_gap = b.Left - (a.Left + a.Width);
    var isInRangeHorizontal = h_gap <= maxHGap;
    var v_gap = Math.abs(a.Top - b.Top);
    var isInRangeVertical = v_gap <= maxVGap;
    return isInRangeHorizontal && isInRangeVertical;
    var _a;
}
exports.isInsameRow = isInsameRow;
function doMerge(a, b) {
    if (a.Left > b.Left) {
        _a = [b, a], a = _a[0], b = _a[1];
    }
    var top = Math.min(a.Top, b.Top);
    var height = Math.max(a.Top + a.Height - top, b.Top + b.Height - top);
    return new dimension_1.Dimension([a.Left, top, (b.Left + b.Width - a.Left), height]);
    var _a;
}
exports.doMerge = doMerge;
function substitudeWithMeasureSpan(doc, n, ctx, start, end) {
    // ignore non-text nodes.
    if (n.nodeType != Node.TEXT_NODE) {
        return;
    }
    removeFollowing(ctx.parent, ctx.index);
    var parent = ctx.parent;
    parent.removeChild(n);
    var preTextNode = doc.createTextNode(n.textContent.substring(0, start));
    var selText = n.textContent.substring(start, end);
    var textNode = doc.createTextNode(selText);
    var postTextNode = doc.createTextNode(n.textContent.substring(end, n.textContent.length));
    parent.appendChild(preTextNode);
    var wrapper = wrapText(doc, textNode);
    parent.appendChild(wrapper);
    parent.appendChild(postTextNode);
    appendFollowing(parent, ctx.index, ctx.siblings);
    return wrapper;
}
exports.substitudeWithMeasureSpan = substitudeWithMeasureSpan;
function restoreBeforeMeasureStatus(parent, nodeIdx, siblings) {
    removeFollowing(parent, nodeIdx - 1);
    appendFollowing(parent, nodeIdx - 1, siblings);
}
exports.restoreBeforeMeasureStatus = restoreBeforeMeasureStatus;
function removeFollowing(parent, startIdx) {
    var children = dom_helper.asArray(parent.childNodes);
    for (var i = children.length - 1; i > startIdx; i--) {
        parent.removeChild(children[i]);
    }
}
exports.removeFollowing = removeFollowing;
function appendFollowing(parent, nodeIdx, siblings) {
    for (var i = nodeIdx + 1; i < siblings.length; i++) {
        parent.appendChild(siblings[i]);
    }
}
exports.appendFollowing = appendFollowing;
function wrapText(doc, textNode) {
    var span = doc.createElement('span');
    var text = textNode.textContent;
    for (var i = 0; i < text.length; i++) {
        var charSpan = doc.createElement('span');
        charSpan.className = constants.MeasureSpanClass();
        var charNode = doc.createTextNode(text.charAt(i));
        charSpan.appendChild(charNode);
        adaptMeasureSpanLayout(charSpan);
        span.appendChild(charSpan);
    }
    adaptMeasureSpanLayout(span);
    return span;
}
exports.wrapText = wrapText;
function getContentOffsets(w, doc, e) {
    var rect = e.getBoundingClientRect();
    var styles = w.getComputedStyle(e);
    var borderLeft = 0;
    var borderTop = 0;
    var boxSizing = dom_helper_1.getBoxSizing(w, e);
    if (boxSizing == dom_helper_1.BoxSizing.ContentBox) {
        borderLeft = dom_helper_1.getStyleNumber(styles, "borderLeftWidth");
        borderTop = dom_helper_1.getStyleNumber(styles, "borderTopWidth");
    }
    var fixed = false;
    var scrollY = w.pageYOffset;
    var scrollX = w.pageXOffset;
    if (styles.position == "fixed") {
        scrollX = 0;
        scrollY = 0;
        fixed = true;
    }
    var dim = new dimension_1.Dimension([
        rect.left + borderLeft + scrollX,
        rect.top + borderTop + scrollY,
        rect.width,
        rect.height
    ]);
    if (fixed) {
        dim.setFixed();
    }
    return dim;
}
function computeLayout(doc, anchor, e) {
    var w = doc.defaultView;
    var eOffset = getContentOffsets(w, doc, e);
    var anchorStyles = w.getComputedStyle(anchor);
    if (anchorStyles.position == "relative" && !eOffset.Fixed) {
        var anchorOffset = getContentOffsets(w, doc, anchor);
        return new dimension_1.Dimension([
            eOffset.Left - anchorOffset.Left,
            eOffset.Top - anchorOffset.Top,
            eOffset.Width,
            eOffset.Height
        ]);
    }
    return eOffset;
}
exports.computeLayout = computeLayout;
function getElementOffsetLayout(doc, el) {
    return new dimension_1.Dimension([el.offsetLeft, el.offsetTop, el.offsetWidth, el.offsetHeight]);
}
exports.getElementOffsetLayout = getElementOffsetLayout;
function adaptMeasureSpanLayout(sp, debug) {
    if (debug === void 0) { debug = true; }
    sp.style.margin = '0';
    sp.style.padding = '0';
    sp.style.border = '0';
    sp.style.font = 'inherit';
    sp.style.fontSize = 'inherit';
    sp.style.verticalAlign = 'inherit';
    if (debug) {
        sp.style.backgroundColor = 'lightgreen';
    }
}
exports.adaptMeasureSpanLayout = adaptMeasureSpanLayout;


/***/ }),
/* 18 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var node_context_1 = __webpack_require__(16);
var range_cache_1 = __webpack_require__(4);
var log_1 = __webpack_require__(1);
var document = null;
/**
 * Iterate through nodes in a range in the following order:
 *   1. first child
 *   2. self
 *   3. next sibling
 *   4. parent's next sibling
 * , until last node in range is met.
 *
 * @param range is a cached object of range.
 *        You should use RangeCache.make(range) to create it.
 * @param visit is a function to be called when a node is visited.
 */
function iterateRangeNodes(range, visit) {
    var endOffset = range.startContainer == range.endContainer ?
        range.endOffset :
        getNodeTextEndIdx(range.startContainer);
    var q = [[range.startContainer, range.startOffset, endOffset]];
    while (q.length > 0) {
        var _a = q.shift(), node = _a[0], start = _a[1], end = _a[2];
        if (node.hasChildNodes()) {
            var next = node.firstChild;
            var endOffset_1 = next == range.endContainer ? range.endOffset :
                getNodeTextEndIdx(next);
            q.push([next, 0, endOffset_1]);
            continue;
        }
        var ctx = new node_context_1.NodeContext(node, range);
        visit(node, ctx, start, end);
        if (node == range.endContainer || node == range.commonAncestorContainer) {
            return;
        }
        if (ctx.nextSibling) {
            var endOffset_2 = ctx.nextSibling == range.endContainer ?
                range.endOffset :
                getNodeTextEndIdx(ctx.nextSibling);
            q.push([ctx.nextSibling, 0, endOffset_2]);
        }
        else {
            tracebackParentNodes(ctx.parent, range, q);
        }
    }
    log_1.error("lastChild", range.commonAncestorContainer.lastChild);
    log_1.error("endContainer", range.endContainer);
    log_1.error("startContainer", range.startContainer);
    log_1.error("commonAncestorContainer:", range.commonAncestorContainer);
    throw 'iterateRangeNodes: end of sub dom tree met.';
}
exports.iterateRangeNodes = iterateRangeNodes;
/**
 * Correct range if needed.
 * For example, make sure the range ends at a leaf node.
 * @param rc
 */
function correctRange(rc) {
    if (rc.endContainer.nodeType == Node.ELEMENT_NODE &&
        rc.endOffset == 0) {
        log_1.debug("range needs correct.");
        var newEnd = rc.endContainer.firstChild;
        var newEndOffset = 0;
        return new range_cache_1.RangeCache(rc.document, rc.commonAncestorContainer, rc.startContainer, newEnd, rc.startOffset, newEndOffset, rc.meta);
    }
    else {
        return rc;
    }
}
exports.correctRange = correctRange;
/**
 * A helper function to determin if the given object
 * is an instance of Range class.
 * @param r an instance of Range of RangeCache
 */
function isRange(r) {
    var asRange = r;
    return asRange.collapse !== undefined && asRange.setStart !== undefined;
}
exports.isRange = isRange;
function getNodeTextEndIdx(node) {
    return node.textContent.length;
}
function tracebackParentNodes(parent, range, q) {
    while (parent != range.commonAncestorContainer) {
        if (parent.nextSibling) {
            var next = parent.nextSibling;
            var endOffset = next == range.endContainer ? range.endOffset :
                next.textContent.length;
            q.push([next, 0, endOffset]);
            break;
        }
        else {
            parent = parent.parentNode;
        }
    }
}


/***/ }),
/* 19 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
/**
 * RangeMeta stores UPath informations needed to store and restore a block.
 * UPath is a string used to uniquely locate a node in the dom tree.
 * It's grammar is defined as:
 *    UPath = [<empty>|#<id>]([/<nodeIndex>])+
 */
var RangeMeta = (function () {
    function RangeMeta(_a, _b) {
        var anchorUPath = _a[0], startNodeUPath = _a[1], endNodeUPath = _a[2], text = _a[3];
        var startCharIndex = _b[0], endCharIndex = _b[1];
        _c = [anchorUPath, startNodeUPath, endNodeUPath, text], this.m_anchorUPath = _c[0], this.m_startNodeUPath = _c[1], this.m_endNodeUPath = _c[2], this.m_text = _c[3];
        _d = [startCharIndex, endCharIndex], this.m_startCharIndex = _d[0], this.m_endCharIndex = _d[1];
        var _c, _d;
    }
    Object.defineProperty(RangeMeta.prototype, "anchorUPath", {
        get: function () {
            return this.m_anchorUPath;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeMeta.prototype, "startNodeUPath", {
        get: function () {
            return this.m_startNodeUPath;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeMeta.prototype, "endNodeUPath", {
        get: function () {
            return this.m_endNodeUPath;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeMeta.prototype, "startCharIndex", {
        get: function () {
            return this.m_startCharIndex;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeMeta.prototype, "endCharIndex", {
        get: function () {
            return this.m_endCharIndex;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(RangeMeta.prototype, "text", {
        get: function () {
            return this.m_text;
        },
        enumerable: true,
        configurable: true
    });
    return RangeMeta;
}());
exports.RangeMeta = RangeMeta;


/***/ }),
/* 20 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
function extend(a, b) {
    for (var _i = 0, b_1 = b; _i < b_1.length; _i++) {
        var o = b_1[_i];
        a.push(o);
    }
}
exports.extend = extend;


/***/ }),
/* 21 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
function generateRandomUUID() {
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() +
        s4() + s4();
}
exports.generateRandomUUID = generateRandomUUID;
function s4() {
    return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
}


/***/ }),
/* 22 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var Markdown = __webpack_require__(11);
var decorator_1 = __webpack_require__(10);
var log_1 = __webpack_require__(1);
var selector_1 = __webpack_require__(8);
var React = __webpack_require__(2);
var ReactDOM = __webpack_require__(12);
var Footnote_1 = __webpack_require__(7);
var decorated_block_1 = __webpack_require__(9);
var dom_helper_1 = __webpack_require__(0);
/**
 * A demo app generate highlights of selected text.
 */
var App = (function () {
    function App(doc) {
        this.m_ftnotes = {};
        this.m_doc = doc;
        this.m_win = doc.defaultView;
    }
    /**
     * If user selected some text in the current window, highlight it.
     * Else do nothing.
     * @return the generated Block id or null.
     */
    App.prototype.highlightSelection = function () {
        var block = selector_1.extractSelectedBlock(window, document);
        if (block == null) {
            return null;
        }
        var decoratedBlock = this.decorate(block);
        this.addBlock(decoratedBlock);
        if (this.m_onHighlight) {
            this.m_onHighlight(decoratedBlock);
        }
        return decoratedBlock.id;
    };
    /**
     * @param handler called when a new DecoratedBlock object is built.
     * It is not guarenteed to be called when the block is rendered in DOM.
     */
    App.prototype.setOnHighlightHandler = function (handler) {
        this.m_onHighlight = handler;
    };
    /**
     * Rebuild the block with a given meta
     * @param meta See @class{RangeMeta}
     * @param id optional, the id for the restored block
     * @return restored block id
     */
    App.prototype.restoreHighlight = function (meta, id) {
        var block = selector_1.restoreBlock(this.m_win, this.m_doc, meta, id);
        var decoratedBlock = this.decorate(block);
        this.addBlock(decoratedBlock);
        return decoratedBlock.id;
    };
    /**
     * Find and remove the block with a given id.
     * @param id of the block to remove
     */
    App.prototype.removeHighlight = function (id) {
        for (var ftId in this.m_ftnotes) {
            var note = this.m_ftnotes[ftId];
            if (note.removeHighlight(id)) {
                return;
            }
        }
        log_1.warn("removeHighlight: " + id + " not found in any footnotes");
    };
    App.prototype.removeAll = function () {
        for (var ftId in this.m_ftnotes) {
            var note = this.m_ftnotes[ftId];
            note.removeAll();
        }
    };
    App.prototype.getAllText = function () {
        var notes = [];
        for (var key in this.m_ftnotes) {
            var note = this.m_ftnotes[key];
            var blocks = note.getAll();
            for (var _i = 0, blocks_1 = blocks; _i < blocks_1.length; _i++) {
                var blk = blocks_1[_i];
                notes.push(blk.text);
            }
            notes.push("");
        }
        return notes;
    };
    App.prototype.generateMarkdownNotes = function () {
        var title = this.m_doc.title || this.m_doc.head.title;
        var url = this.m_win.location.toString();
        var lines = this.getAllText();
        var noteTitle = "Notes on " + Markdown.link(title, url) + "\r\n\r\n";
        var md = Markdown.h(1, noteTitle);
        var buf = [];
        for (var idx = 0; idx <= lines.length; idx++) {
            if (idx == lines.length || lines[idx] == "") {
                if (buf.length > 0) {
                    md += Markdown.ul(buf, 0);
                }
                md += "\r\n";
                buf.length = 0;
            }
            else {
                buf.push(lines[idx]);
            }
        }
        return md;
    };
    App.prototype.decorate = function (block) {
        var factory = new decorated_block_1.DecoratedBlockFactory(block);
        factory.addDecorator(decorator_1.BasicDecorator);
        return factory.make();
    };
    App.prototype.addBlock = function (blk) {
        var ftntId = blk.rangeCache.meta.anchorUPath;
        if (ftntId in this.m_ftnotes) {
            var ftnt = this.m_ftnotes[ftntId];
            ftnt.addHighlight(blk);
            log_1.debug("added to existing ftnt: " + ftntId);
        }
        else {
            this.renderFootnote(ftntId, [blk]);
        }
    };
    App.prototype.getOrCreateContainer = function (cid, parent) {
        var d = this.m_doc.getElementById(cid);
        if (!d) {
            d = this.m_doc.createElement("div");
            d.id = cid;
            parent.appendChild(d);
        }
        return d;
    };
    App.prototype.renderFootnote = function (anchorUPath, initBlocks) {
        var _this = this;
        log_1.debug("generating new ftnt under : " + anchorUPath);
        var anchor = dom_helper_1.getNodeByPath(this.m_doc, anchorUPath);
        var containerId = "ftnt_" + anchorUPath;
        var container = this.getOrCreateContainer(containerId, anchor);
        var ftntEl = React.createElement(Footnote_1.Footnote, {
            id: anchorUPath,
            source: { get: function () { return initBlocks; }, poll: false },
            onMount: function (ins) {
                _this.m_ftnotes[anchorUPath] = ins;
            }
        });
        log_1.debug("generated under", container);
        ReactDOM.render(ftntEl, container);
    };
    return App;
}());
exports.App = App;


/***/ })
/******/ ]);
//# sourceMappingURL=bundle.js.map