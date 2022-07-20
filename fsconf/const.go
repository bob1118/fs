package fsconf

const NOT_FOUND = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="result">
    <result status="not found"/>
  </section>
</document>`

///////////////////////////////////////////////////////////////////////////////
//////////////////////////////section: configuration///////////////////////////
///////////////////////////////////////////////////////////////////////////////
const CONFIGURATION = `<document type="freeswitch/xml"  encoding="UTF-8">
<section name="configuration">
%s
</section>
</document>`

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////section: dialplan////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// esl inbound default action park.
// profile internal context=public => directory useragent context=default
// profile external context=public
const DIALPLAN_APP_PARK = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="dialplan" description="dialplan inbound for FreeSwitch">
    <context name="default">
      <extension name="default">
        <condition>
          <action application="set" data="continue_on_fail=true"/>
          <action application="park"/>
        </condition>
      </extension>
    </context>
    <context name="public">
      <extension name="default">
        <condition>
          <action application="set" data="continue_on_fail=true"/>
          <action application="park"/>
        </condition>
      </extension>
    </context>
  </section>
</document>
`

// esl outbound default action socket.
// profile internal context=public => directory useragent context=default
// profile external context=public
const DIALPLAN_APP_SOCKET = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="dialplan" description="dialplan outbound FreeSwitch">
    <context name="default">
      <extension name="default">
        <condition>
          <action application="socket" data="%s async full"/>
        </condition>
      </extension>
    </context>
    <context name="public">
      <extension name="default">
        <condition>
          <action application="socket" data="%s async full"/>
        </condition>
      </extension>
    </context>
  </section>
</document>
`

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////section: direcotry////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
