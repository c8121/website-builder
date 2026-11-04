import {defineAsyncComponent, nextTick} from 'vue'

export default {

    components: {
        TextEditor: defineAsyncComponent(() => import('./content-editor/text.vue.js')),
        TestEditor: defineAsyncComponent(() => import('./content-editor/test.vue.js'))
    },

    template: `
        <div>
            <h1>PageView</h1>
            <pre>{{ page }}</pre>
            
            <div v-if="page">

                <div class="mb-2">
                    <page-select v-model="page.ParentId" ref="pageSelect" label="Parent" @component-event="treeEvent" />
                </div>

                <div class="input-group mb-2">
                    <label class="input-group-text">Name</label>
                    <input v-model="page.Name" class="form-control">
                </div>
                
                <div>
                    <select v-model="selectedEditor" class="form-select">
                        <option :value="null"></option>
                        <option v-for="editor in availableEditors" :value="editor">{{ editor.name }}</option>
                    </select>
                </div>
                
                <div class="mt-4">
                    <component :is="selectedEditor ? selectedEditor.component : null" ref="editor" @vue:mounted="editorReady"></component>
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
            page: null,
            availableEditors: [
                {name: 'Text', component: 'TextEditor'},
                {name: 'Test', component: 'TestEditor'}
            ],
            selectedEditor: null
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
                        self.loadEditor();
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
                Name: "New Page",
                ContentType: "Text",
                Content: ""
            }

            self.loadEditor();

            if (self.$refs.pageSelect)
                self.$refs.pageSelect.reloadPageTree();
        },

        loadEditor() {
            const self = this;

            if (!self.page.ContentType) {
                console.log("No ContentType given");
                return;
            }

            self.selectedEditor = self.getEditorByName(self.page.ContentType);
        },

        editorReady() {
            this.$refs.editor.setContent(this.page.Content);
        },

        getEditorByName(name) {
            for (const editor of this.availableEditors) {
                if (name === editor.name)
                    return editor;
            }
            console.log("No such editor: " + name);
            return null;
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

        treeEvent() {
            console.log(arguments);
        }
    },

    emits: ['componentEvent'],
}