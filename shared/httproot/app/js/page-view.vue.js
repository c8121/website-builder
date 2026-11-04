export default {
    template: `
        <div>
            <h1>PageView</h1>
            <pre>{{ page }}</pre>

            <div v-if="page">

                <div class="mb-2">
                    <page-select v-model="page.ParentId" label="Parent" @component-event="treeEvent" />
                </div>

                <div class="input-group">
                    <label class="input-group-text">Name</label>
                    <input v-model="page.Name" class="form-control">
                </div>

                <div class="mt-4">
                    <button class="btn btn-primary" @click="savePage">Save</button>
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

            if(self.parentId) {
                self.createPage();
                return;
            }

            if(!self.id)
                return;

            const requestOptions = {
                method: 'GET'
            }
            fetch('/page/' + self.id, requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.page = json;
                });
        },

        createPage() {
            const self = this;
            if(!self.parentId)
                return;

            self.page = {
                ParentId: parseInt(self.parentId),
                Name: "New Page"
            }
        },

        savePage() {
            const self = this;
            if(!self.page)
                return;

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(self.page)
            }
            fetch('/page/add', requestParams)
                .then(res => res.json())
                .then(json => {
                    if( json.Id) {
                        self.page = json
                        this.$emit('componentEvent', 'pageSaved', 'page-view', self.page);
                    } else if(json.Message) {
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