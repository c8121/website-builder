export default {
    template: `
        <div>
            <div class="dropdown">
                <div class="input-group">
                    <label v-if="label" class="input-group-text">{{label}}</label>
                    <button class="form-control" data-bs-toggle="dropdown" data-bs-auto-close="outside" aria-expanded="false">
                        <i class="bi bi-folder-fill"></i><span v-if="page">{{page.Name}}</span>
                    </button>
                    <div class="dropdown-menu ps-3">
                        <page-tree ref="pageTree" @component-event="componentEvent" />
                    </div>
                    <button v-if="page" @click="clearValue" class="btn border"><i class="bi bi-x-circle"></i></button>
                </div>
            </div>
        </div>
    `,

    props: {
        modelValue: {
            type: Number,
            default: null
        },
        label: {
            type: String,
            default: null
        }
    },

    watch: {
        modelValue: {
            handler(newValue) {
                if(!this.page || this.page.Id !== newValue)
                    this.loadPage(newValue);
            },
            immediate: true
        }
    },

    data() {
        return {
            page: null
        }
    },

    methods: {

        loadPage(id) {
            const self = this;
            if(!id) {
                self.page = null;
                return;
            }

            const requestOptions = {
                method: 'GET'
            }
            fetch('/page/' + id, requestOptions)
                .then(res => res.json())
                .then(json => {
                    if(json.Id) {
                        self.page = json;
                    } else {
                        this.$emit('update:modelValue', 0);
                    }
                    console.log(json);
                });
        },

        clearValue() {
            this.$emit('update:modelValue', 0);
        },

        componentEvent(eventName, sender) {
            const self = this;
            switch (eventName) {
                case 'pageClick': {
                    const page = arguments[2];
                    this.$emit('update:modelValue', page.Id);
                    break;
                }
            }
        }
    },

    emits: ['componentEvent', 'update:modelValue'],
}