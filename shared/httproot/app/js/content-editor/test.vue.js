export default {
    template: `
        <div>
            <pre>{{editorContent}}</pre>
            
            <div v-if="editorContent && editorContent.text">
                <input v-model="editorContent.text">
            </div>
            
            <button @click="editorContent = { text: 'foobar'}">{}</button>
            <button @click="editorContent = [];">[]</button>
            <button @click="addItem">+</button>
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
                    this.editorContent = JSON.parse(content)
                } catch (e) {
                    console.log("Failed to parse content");
                    this.editorContent = {Content: content}
                }
            } else {
                this.editorContent = {}
            }
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