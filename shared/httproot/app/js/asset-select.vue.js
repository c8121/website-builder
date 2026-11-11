import { selectAssetDialog } from './asset-select-dialog.js'

export default {
    template: `
        <div>
            <div v-if="url && previewUrl()" class="form-control p-2 mb-1 overflow-scroll">
                <img :src="previewUrl()">
            </div>
            <div class="input-group">
                <label v-if="label" class="input-group-text">{{label}}</label>
                <input class="form-control" v-model="url" @input="urlChanged">
                <button class="btn border" @click="selectAsset"><i class="bi bi-image"></i></button>
            </div>
        </div>
    `,

    props: {
        modelValue: {
            type: String,
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
                this.url = newValue;
            },
            immediate: true
        }
    },

    data() {
        return {
            url: ""
        }
    },

    methods: {

        selectAsset() {
            selectAssetDialog("Select asset", (url) => {
                this.$emit('update:modelValue', url);
            });
        },

        urlChanged() {
            this.$emit('update:modelValue', this.url);
        },

        previewUrl() {
            if (this.url && this.url.indexOf("/assets/") === 0) {
                const p = this.url.lastIndexOf("/");
                if (p > -1) {
                    const hash = this.url.substring(p + 1)
                    if (hash.length > 30)
                        return "/assets/thumbnail/" + hash;
                }
            }
            return "";
        }
    },

    emits: ['update:modelValue'],
}