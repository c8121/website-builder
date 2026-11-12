import {defineAsyncComponent, nextTick} from 'vue'

export default {

    components: {
        HtmlEditor: defineAsyncComponent(() => import('./content-editor/html.vue.js')),
        MarkdownEditor: defineAsyncComponent(() => import('./content-editor/md.vue.js')),
        TextEditor: defineAsyncComponent(() => import('./content-editor/text.vue.js')),
        ImageEditor: defineAsyncComponent(() => import('./content-editor/image.vue.js')),
        TestEditor: defineAsyncComponent(() => import('./content-editor/test.vue.js'))
    },

    template: `
        <div>
            <div v-if="content">

                <div v-if="content.PageId" class="mb-3">
                    <a class="btn btn-sm btn-link" :href="'#/page/' + content.PageId"><i class="bi bi-arrow-left-circle"></i> back</a>
                </div>

                <div class="mb-2">
                    <page-select v-model="content.PageId" ref="pageSelect" label="Page" @component-event="treeEvent" />
                </div>

                <div class="input-group mb-4">
                    <label class="input-group-text">Name</label>
                    <input v-model="content.Name" class="form-control">
                </div>
                
                <div>
                    <select v-model="selectedEditor" class="form-select">
                        <option :value="null"></option>
                        <option v-for="editor in availableEditors" :value="editor">{{ editor.caption }}</option>
                    </select>
                </div>
                
                <div class="mt-4">
                    <component :is="selectedEditor ? selectedEditor.component : null" ref="editor" @vue:mounted="editorReady"></component>
                </div>

                <div class="mt-4">
                    <button class="btn btn-primary" @click="savePageContent">Save</button>
                    <button v-if="content.Id" class="btn btn-outline-danger ms-2" @click="deletePageContent">Delete</button>
                </div>
            </div>

        </div>
    `,

    props: {
        id: {
            type: String,
            default: null
        },
        pageId: {
            type: String,
            default: null
        }
    },

    data() {
        return {
            content: null,
            availableEditors: [
                {caption: 'Rich-Text', name: 'Html', component: 'HtmlEditor'},
                {caption: 'Markdown', name: 'Markdown', component: 'MarkdownEditor'},
                {caption: 'Text', name: 'Text', component: 'TextEditor'},
                {caption: 'Image', name: 'Image', component: 'ImageEditor'},
                {caption: 'Test', name: 'Test', component: 'TestEditor'}
            ],
            selectedEditor: null
        }
    },

    watch: {
        id() {
            this.loadPageContent();
        },
        pageId() {
            this.createPageContent();
        }
    },

    mounted() {
        this.loadPageContent();
    },

    methods: {

        loadPageContent() {
            const self = this;

            if (self.pageId) {
                self.createPageContent();
                return;
            }

            if (!self.id)
                return;

            const requestOptions = {
                method: 'GET'
            }
            fetch('/content/' + self.id, requestOptions)
                .then(res => res.json())
                .then(json => {
                    if (json.Id) {
                        self.content = json;
                        self.loadEditor();
                    } else {
                        alert(json)
                    }
                });
        },

        createPageContent() {
            const self = this;
            if (!self.pageId)
                return;

            self.content = {
                PageId: parseInt(self.pageId),
                Name: "New Content",
                ContentType: "Text",
                Content: ""
            }

            self.loadEditor();
        },

        loadEditor() {
            const self = this;

            if (!self.content.ContentType) {
                console.log("No ContentType given");
                return;
            }

            self.selectedEditor = self.getEditorByName(self.content.ContentType);
        },

        editorReady() {
            this.$refs.editor.setContent(this.content.Content);
        },

        getEditorByName(name) {
            for (const editor of this.availableEditors) {
                if (name === editor.name)
                    return editor;
            }
            console.log("No such editor: " + name);
            return null;
        },

        savePageContent() {
            const self = this;
            if (!self.content)
                return;

            if (self.selectedEditor) {
                self.content.ContentType = self.selectedEditor.name;
                self.content.Content = self.$refs.editor.getContent();
            }

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(self.content)
            }
            fetch('/content/add', requestParams)
                .then(res => res.json())
                .then(json => {
                    if (json.Id) {
                        self.content = json
                        this.$emit('componentEvent', 'pageContentSaved', 'page-content-view', self.content);
                    } else if (json.Message) {
                        alert(json.Message)
                    }
                });
        },

        deletePageContent() {
            const self = this;
            if (!self.content)
                return;

            if (!confirm("Delete?"))
                return;

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: "id=" + encodeURI(self.content.Id)
            }
            fetch('/content/delete', requestParams)
                .then(res => res.json())
                .then(json => {
                    if (json.Status && json.Status === 'ok') {
                        self.content = null
                        this.$emit('componentEvent', 'pageDeleted', 'page-content-view', json);
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