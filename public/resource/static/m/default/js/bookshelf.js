const STORAGE_FAVORITE_KEY = "footprint";
const STORAGE_FAVORITE_SIZE = 30;

class LRU {
    constructor(cache = new Map(), max = 10) {
        this.max = max;
        this.cache = cache;
    }

    getCache() {
        return this.cache;
    }

    get(key) {
        let item = this.cache.get(key);
        if (item) {
            // refresh key
            this.cache.delete(key);
            this.cache.set(key, item);
        }
        return item;
    }

    set(key, val) {
        // refresh key
        if (this.cache.has(key)) {
            this.cache.delete(key);
        }
        // evict oldest
        else if (this.cache.size === this.max) {
            this.cache.delete(this.first());
        }
        this.cache.set(key, val);
    }

    remove(key) {
        if (this.cache.has(key)) {
            this.cache.delete(key);
        }
    }

    first() {
        return this.cache.keys().next().value;
    }
}

function saveFootprint(chapterInfo) {
    var cacheMap = new Map(JSON.parse(localStorage.getItem(STORAGE_FAVORITE_KEY)))
    var lru = new LRU(cacheMap, STORAGE_FAVORITE_SIZE)

    lru.set(chapterInfo.id, chapterInfo)
    localStorage.setItem(STORAGE_FAVORITE_KEY, JSON.stringify(Array.from(lru.getCache().entries())))
}

function removeFootprint(novelId) {
    var cacheMap = new Map(JSON.parse(localStorage.getItem(STORAGE_FAVORITE_KEY)))
    var lru = new LRU(cacheMap, STORAGE_FAVORITE_SIZE)

    lru.remove(novelId);
    localStorage.setItem(STORAGE_FAVORITE_KEY, JSON.stringify(Array.from(lru.getCache().entries())))
}

function loadBookshelfToHtml(displayNum = STORAGE_FAVORITE_SIZE, displayAll=true) {
    var cacheMap = new Map(JSON.parse(localStorage.getItem(STORAGE_FAVORITE_KEY)))
    if (cacheMap.size === 0) {
        $("#footprint").css("display", "none")
        return;
    }

    $("#footprint").css("display", "block")
    var html = ""
    var i = 1;
    for (let chapter of Array.from(cacheMap.values()).reverse()) {
        html += "<tr>";
        if (displayAll) {
            html += '<td>' + i + "</td>";
        }
        html += '<td><a href="' + chapter.novUrl + '">' + chapter.novName + "</a></td>";
        html += '<td><a href="' + chapter.chapterUrl + '">' + chapter.chapterTitle + "</a></td>";
        html += '<td>' + chapter.author + "</td>";
        if (displayAll) {
            html += '<td class="bookshelf-action"><a class="pointer text-red" href="javascript:;" onclick="removeFromBookshelf(this, ' + chapter.id + ')">Remove</a></td>';
        }
        html += "</tr>";
        i++;
        if (i > displayNum) break;
    }

    $("#footprint").find("tbody").append(html)
}

function removeFromBookshelf(el, id) {
    var td = el.closest("tr");
    td.remove();
    removeFootprint(id);
}