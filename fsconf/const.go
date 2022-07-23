package fsconf

const NOT_FOUND = `
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="result">
    <result status="not found"/>
  </section>
</document>`

///////////////////////////////////////////////////////////////////////////////
//////////////////////////////section: configuration///////////////////////////
///////////////////////////////////////////////////////////////////////////////

const CONFIGURATION = `
<document type="freeswitch/xml" encoding="UTF-8">
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
</document>`

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
</document>`

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////section: direcotry////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

const DIRECTORY = `
<document type="freeswitch/xml" encoding="UTF-8">
  <section name="directory">
    %s
  </section>
</document>`

// const DOMAIN = `
//     <domain name="%s">
//       <params>
//         <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"/>
//         <param name="jsonrpc-allowed-methods" value="verto"/>
//       </params>
//       <variables>
//         <variable name="record_stereo" value="true"/>
//       </variables>
//       <groups>
//         <group name="default">
//           <users>
//             <user id="default">
//               <params>
//                 <param name="password" value="$${default_password}"/>
//               </params>
//               <variables>
//                 <variable name="numbering_plan" value="$${default_country}"/>
//                 <variable name="default_areacode" value="$${default_areacode}"/>
//               </variables>
//             </user>
//           </users>
//         </group>
//       </groups>
//     </domain>
// `

const DOMAIN = `
    <domain name="%s">
      <params>
        <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"/>
        <param name="jsonrpc-allowed-methods" value="verto"/>
      </params>
      <groups>
        <group name="default">
          <users>
            <user id="default"/>
          </users>
        </group>
      </groups>
    </domain>
`

const USERAGENT = `<document type="freeswitch/xml" encoding="UTF-8">
<section name="directory">
 <domain name="%s">
  <groups>
  <group name="%s">
   <users>
   <user id="%s"  cacheable="%s">
    <params>
     <param name="password" value="%s"/>
     <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"/>
    </params>
    <variables>
     <variable name="user_context" value="default"/>
     <variable name="record_stereo" value="true"/>
    </variables>
   </user>
   </users>
  </group>
  </groups>
 </domain>
</section>
</document>
`

const USERAGENT_A1HASH = `<document type="freeswitch/xml" encoding="UTF-8">
<section name="directory">
 <domain name="%s">
  <groups>
  <group name="%s">
   <users>
   <user id="%s"  cacheable="%s">
    <params>
     <param name="a1-hash" value="%s"/>
     <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"/>
    </params>
    <variables>
     <variable name="user_context" value="default"/>
     <variable name="record_stereo" value="true"/>
    </variables>
   </user>
   </users>
  </group>
  </groups>
 </domain>
</section>
</document>
`

const USERAGENT_REVERSE = `<document type="freeswitch/xml" encoding="UTF-8">
<section name="directory">
 <domain name="%s">
  <groups>
  <group name="%s">
   <users>
   <user id="%s"  cacheable="%s">
    <params>
     <param name="reverse-auth-user" value="%s"/>
     <param name="reverse-auth-pass" value="%s"/>
    </params>
   </user>
   </users>
  </group>
  </groups>
 </domain>
</section>
</document>
`

///////////////////////////////////vars.xml////////////////////////////////////////

const VARS_NEW_PASSWORD_WITH_IPV4_AND_PGHANDLE = `
  <X-PRE-PROCESS cmd="set" data="default_password=D_e_f_a_u_l_t_P_a_s_s_w_o_r_d"/>
  <X-PRE-PROCESS cmd="set" data="local_ip_v4=%s"/>
  <X-PRE-PROCESS cmd="set" data="pg_handle=%s"/>
  <X-PRE-PROCESS cmd="set" data="json_db_handle=$${pg_handle}"/>
`
