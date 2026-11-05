export default {
    template: `
        <div>
            <template v-if="editorContent">
                <div class="mb-2">
                    <asset-select label="URL" v-model="editorContent.src"></asset-select>
                </div>
                <div class="input-group">
                    <label class="input-group-text">Alt</label>
                    <input v-model="editorContent.alt" class="form-control">
                </div>
            </template>
        </div>
    `,

    data() {
        return {
            editorContent: null
        }
    },


    methods: {
        setContent(content) {
            if (content) {
                try {
                    this.editorContent = this.initContent(JSON.parse(content));
                } catch (e) {
                    console.log("Failed to parse content");
                    this.editorContent = this.initContent({Content: content});
                }
            } else {
                this.editorContent = this.initContent({});
            }
        },

        initContent(content) {
            if(!content.src) {
                content.src = "";
            }
            if(!content.alt) {
                content.alt = "";
            }

            return content;
        },

        getContent() {
            return JSON.stringify(this.editorContent);
        },

        addItem() {
            var val = this.editorContent;
            val.push("TEST")
            this.$emit('update:modelValue', {
                ...this.modelValue,
                Content: JSON.stringify(val)
            });
        }
    },

    emits: ['componentEvent'],
}