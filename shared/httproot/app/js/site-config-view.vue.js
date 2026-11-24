import {defineAsyncComponent, nextTick} from 'vue'
import {EditState} from './edit-state.vue.js'

export default {

    mixins: [EditState],

    components: {
        HtmlEditor: defineAsyncComponent(() => import('./config-editor/html.vue.js')),
        TextEditor: defineAsyncComponent(() => import('./config-editor/text.vue.js')),
        ImageEditor: defineAsyncComponent(() => import('./config-editor/image.vue.js'))
    },

    template: `
        <div>
            <div v-if="config">

                <div class="input-group mb-4">
                    <label class="input-group-text">Name</label>
                    <input v-model="config.Name" class="form-control">
                </div>
                
                
                <div class="input-group mb-4">
                    <label class="input-group-text">Type</label>
                    <select v-model="selectedEditor" class="form-select">
                        <option :value="null"></option>
                        <option v-for="editor in availableEditors" :value="editor">{{ editor.caption }}</option>
                    </select>
                </div>
                
                <div class="mt-4">
                    <component :is="selectedEditor ? selectedEditor.component : null" ref="editor" @vue:mounted="editorReady"></component>
                </div>

                <div class="mt-4">
                    <button class="btn btn-primary" @click="saveConfig">Save</button>
                    <button v-if="config.Id" class="btn btn-outline-danger ms-2" @click="deleteConfig">Delete</button>
                </div>
            </div>

        </div>
    `,

    props: {
        id: {
            type: String,
            default: null
        }
    },

    data() {
        return {
            config: null,
            availableEditors: [
                {caption: 'Rich-Text', name: 'Html', component: 'HtmlEditor'},
                {caption: 'Text', name: 'Text', component: 'TextEditor'},
                {caption: 'Image', name: 'Image', component: 'ImageEditor'}
            ],
            selectedEditor: null
        }
    },

    watch: {
        id() {
            this.loadSiteConfig();
        }
    },

    beforeRouteLeave() {
        return this.confirmLeave();
    },
    beforeRouteUpdate() {
        return this.confirmLeave();
    },

    mounted() {
        this.loadSiteConfig();
    },

    methods: {

        getState() {
            if (!this.config)
                return null;
            this.updateConfig();
            return this.config;
        },

        loadSiteConfig() {
            const self = this;

            if (!self.id) {
                self.createSiteConfig();
                return;
            }

            const requestOptions = {
                method: 'GET'
            }
            fetch('/siteconfig/' + self.id, requestOptions)
                .then(res => res.json())
                .then(json => {
                    if (json.Id) {
                        self.config = json;
                        self.saveState();
                        self.loadEditor();
                    } else {
                        alert(json)
                    }
                });
        },

        createSiteConfig() {
            const self = this;
            
            self.config = {
                Name: "",
                ContentType: "Text",
                Value: ""
            }

            self.loadEditor();
        },

        loadEditor() {
            const self = this;

            if (!self.config.ContentType) {
                console.log("No ContentType given");
                self.config.ContentType = "Text";
            }

            self.selectedEditor = self.getEditorByName(self.config.ContentType);
        },

        editorReady() {
            this.$refs.editor.setValue(this.config.Value);
        },

        getEditorByName(name) {
            for (const editor of this.availableEditors) {
                if (name === editor.name)
                    return editor;
            }
            console.log("No such editor: " + name);
            return null;
        },

        updateConfig() {
            if (this.selectedEditor) {
                this.config.ContentType = this.selectedEditor.name;
                this.config.Value = this.$refs.editor.getValue();
            }
        },

        saveConfig() {
            const self = this;
            if (!self.config)
                return;
            self.updateConfig();

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(self.config)
            }
            fetch('/siteconfig/add', requestParams)
                .then(res => res.json())
                .then(json => {
                    if (json.Id) {
                        self.config = json
                        self.saveState();
                        this.$emit('componentEvent', 'siteConfigSaved', 'site-config-view', self.config);
                    } else if (json.Message) {
                        alert(json.Message)
                    }
                });
        },

        deleteConfig() {
            const self = this;
            if (!self.config)
                return;

            if (!confirm("Delete?"))
                return;

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: "id=" + encodeURI(self.config.Id)
            }
            fetch('/siteconfig/delete', requestParams)
                .then(res => res.json())
                .then(json => {
                    if (json.Status && json.Status === 'ok') {
                        self.config = null
                        this.$emit('componentEvent', 'configDeleted', 'site-config-view', json);
                    } else if (json.Message) {
                        alert(json.Message)
                    }
                });
        },
    },

    emits: ['componentEvent'],
}