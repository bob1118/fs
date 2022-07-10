package fsconf

const NOT_FOUND = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="result">
    <result status="not found"/>
  </section>
</document>`

const CONFIGURATION = `<document type="freeswitch/xml"  encoding="UTF-8">
<section name="configuration">
%s
</section>
</document>`
