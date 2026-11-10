export default {


    template: `
        <div>
            <div v-if="page">

                <div class="mb-2">
                    <page-select v-model="page.ParentId" ref="pageSelect" label="Parent" @component-event="treeEvent" />
                </div>

                <div class="input-group mb-2">
                    <label class="input-group-text">Name</label>
                    <input v-model="page.Name" class="form-control">
                </div>
                
                <div class="card" v-if="page.Id">
                    <div class="card-header">Content</div>
                    <div class="card-body">
                        <page-content-list :value="page.Id" ref="contentList" @component-event="listEvent" />
                    </div>
                    <div class="card-footer">
                        <button class="btn btn-sm btn-secondary" @click="addContent"><i class="bi bi-file-earmark-plus"></i> Add</button>
                    </div>
                </div>

                <div class="mt-4">
                    <button class="btn btn-primary" @click="savePage">Save</button>
                    <button v-if="page.Id" class="btn btn-outline-danger ms-2" @click="deletePage">Delete</button>
                </div>
            </div>

        </div>
    `,

    props: {
        id: {
            type: String,
            default: null
        },
        parentId: {
            type: String,
            default: null
        }
    },

    data() {
        return {
            page: null
        }
    },

    watch: {
        id() {
            this.loadPage();
        },
        parentId() {
            this.createPage();
        }
    },

    mounted() {
        this.loadPage();
    },

    methods: {

        loadPage() {
            const self = this;

            if (self.parentId) {
                self.createPage();
                return;
            }

            if (!self.id)
                return;

            const requestOptions = {
                method: 'GET'
            }
            fetch('/page/' + self.id, requestOptions)
                .then(res => res.json())
                .then(json => {
                    if (json.Id) {
                        self.page = json;
                    } else {
                        alert(json)
                    }
                    if (self.$refs.pageSelect)
                        self.$refs.pageSelect.reloadPageTree();
                });
        },

        createPage() {
            const self = this;
            if (!self.parentId)
                return;

            self.page = {
                ParentId: parseInt(self.parentId),
                Name: "New Page"
            }

            if (self.$refs.pageSelect)
                self.$refs.pageSelect.reloadPageTree();
        },

        savePage() {
            const self = this;
            if (!self.page)
                return;

            if (self.selectedEditor) {
                self.page.ContentType = self.selectedEditor.name;
                self.page.Content = self.$refs.editor.getContent();
            }

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(self.page)
            }
            fetch('/page/add', requestParams)
                .then(res => res.json())
                .then(json => {
                    if (json.Id) {
                        self.page = json
                        this.$emit('componentEvent', 'pageSaved', 'page-view', self.page);
                    } else if (json.Message) {
                        alert(json.Message)
                    }
                });
        },

        deletePage() {
            const self = this;
            if (!self.page)
                return;

            if (!confirm("Delete?"))
                return;

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: "id=" + encodeURI(self.page.Id)
            }
            fetch('/page/delete', requestParams)
                .then(res => res.json())
                .then(json => {
                    if (json.Status && json.Status === 'ok') {
                        self.page = null
                        this.$emit('componentEvent', 'pageDeleted', 'page-view', json);
                    } else if (json.Message) {
                        alert(json.Message)
                    }
                });
        },

        addContent() {
            this.$emit('componentEvent', 'addContent', 'page-view', this.page);
        },

        treeEvent() {
            this.$emit('componentEvent', arguments[0], arguments[1], arguments[2]);
        },

        listEvent() {
            this.$emit('componentEvent', arguments[0], arguments[1], arguments[2]);
        },

    },

    emits: ['componentEvent'],
}