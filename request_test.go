package openrtb

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseRequestBytes_Blank(t *testing.T) {
	req, err := ParseRequestBytes([]byte("{}"))
	assert.Nil(t, err)
	assert.IsType(t, &Request{}, req)
}

func TestParseRequestBytes_SimpleBanner(t *testing.T) {
	req, err := ParseRequestBytes(simpleBanner)
	assert.Nil(t, err)
	assert.IsType(t, &Request{}, req)

	assert.Equal(t, *req.At, 2)
	assert.Equal(t, *req.Id, "1234534625254")
	assert.Equal(t, len(req.Badv), 2)
	assert.Equal(t, len(req.Bcat), 0)

	assert.Equal(t, len(req.Imp), 1)
	assert.Equal(t, *req.Imp[0].Banner.W, 300)
	assert.Equal(t, *req.Imp[0].Banner.H, 250)

	assert.Equal(t, *req.Site.Name, "Site ABCD")
	assert.Equal(t, *req.Site.Publisher.Name, "Publisher A")
	assert.Equal(t, *req.Device.Ip, "64.124.253.1")
	assert.Equal(t, *req.User.Buyeruid, "5df678asd8987656asdf78987654")
}

func TestParseRequestBytes_ExpandableCreative(t *testing.T) {
	req, err := ParseRequestBytes(expandableCreative)
	assert.Nil(t, err)
	assert.IsType(t, &Request{}, req)

	assert.Equal(t, *req.At, 2)
	assert.Equal(t, *req.Tmax, 120)
	assert.Equal(t, req.Imp[0].Banner.Expdir, []int{2, 4})
	assert.Equal(t, *req.Site.Privacypolicy, 1)
	assert.Equal(t, *req.Device.Flashver, "10.1")
	assert.Equal(t, len(req.User.Data), 1)
	assert.Equal(t, *req.User.Data[0].Id, "6")
	assert.Equal(t, len(req.User.Data[0].Segment), 3)
	assert.Equal(t, *req.User.Data[0].Segment[2].Id, "23423424")
}

func TestParseRequest_ExpandableCreative(t *testing.T) {
	req, err := ParseRequest(bytes.NewBuffer(expandableCreative))
	assert.Nil(t, err)
	assert.IsType(t, &Request{}, req)
}

func TestRequest_Valid(t *testing.T) {
	r := &Request{}
	s := &Site{}
	a := &App{}
	i := &Impression{}
	b := &Banner{}

	// blank Request
	ok, err := r.Valid()
	assert.Equal(t, ok, false)
	if err != nil {
		assert.Equal(t, err.Error(), "openrtb parse: request ID missing")
	}

	// with ID
	r.SetId("RAND_ID")
	ok, err = r.Valid()
	assert.Equal(t, ok, false)
	if err != nil {
		assert.Equal(t, err.Error(), "openrtb parse: no impressions")
	}

	// with Site
	r.SetSite(*s)
	ok, err = r.Valid()
	assert.Equal(t, ok, false)
	if err != nil {
		assert.Equal(t, err.Error(), "openrtb parse: no impressions")
	}

	// with Site & App
	r.SetApp(*a)
	ok, err = r.Valid()
	assert.Equal(t, ok, false)
	if err != nil {
		assert.Equal(t, err.Error(), "openrtb parse: no impressions")
	}

	// with Impression
	i.SetId("IMPID").SetBanner(*b).WithDefaults()
	r.Imp = []Impression{*i}
	ok, err = r.Valid()
	assert.Equal(t, ok, false)
	if err != nil {
		assert.Equal(t, err.Error(), "openrtb parse: request has site and app")
	}

	// with valid attrs
	r.App = nil
	ok, err = r.Valid()
	assert.Equal(t, ok, true)
}

func TestRequest_WithDefaults(t *testing.T) {
	s := &Site{}
	a := &App{}
	d := &Device{}
	i := &Impression{}
	b := &Banner{}
	v := &Video{}

	i.SetBanner(*b).SetVideo(*v)
	r := &Request{Site: s, App: a, Device: d, Imp: []Impression{*i}}

	req := r.WithDefaults()
	assert.Equal(t, *req.At, 2)
	assert.Equal(t, *req.App.Privacypolicy, 0)
	assert.Equal(t, *req.App.Paid, 0)
	assert.Equal(t, *req.Site.Privacypolicy, 0)
	assert.Equal(t, *req.Device.Dnt, 0)
	assert.Equal(t, *req.Device.Js, 0)
	assert.Equal(t, *req.Device.Connectiontype, CONN_TYPE_UNKNOWN)
	assert.Equal(t, *req.Device.Devicetype, DEVICE_TYPE_UNKNOWN)
	assert.Equal(t, *req.Imp[0].Instl, 0)
	assert.Equal(t, *req.Imp[0].Bidfloor, 0)
	assert.Equal(t, *req.Imp[0].Bidfloorcur, "USD")
	assert.Equal(t, *req.Imp[0].Banner.Topframe, 0)
	assert.Equal(t, *req.Imp[0].Banner.Pos, AD_POS_UNKNOWN)
	assert.Equal(t, *req.Imp[0].Video.Sequence, 1)
	assert.Equal(t, *req.Imp[0].Video.Boxingallowed, 1)
	assert.Equal(t, *req.Imp[0].Video.Pos, AD_POS_UNKNOWN)
}

func TestRequest_JSON(t *testing.T) {
	req, err := ParseRequest(bytes.NewBuffer(expandableCreative))
	assert.Nil(t, err)

	if req != nil {
		assert.Equal(t, "pending", "TODO")
		// json, err := req.JSON()
		// assert.Nil(t, err)
		// assert.Equal(t, string(json), string(expandableCreative))
	}
}

var simpleBanner []byte = []byte(`
{
  "id":"1234534625254",
  "at":2,
  "tmax":120,
  "imp":[
    {
      "id":"1",
      "banner":{
        "w":300,
        "h":250,
        "pos":1,
        "battr":[13]
      }
    }
  ],
  "badv":["company1.com","company2.com"],
  "site":{
    "id":"234563",
    "name":"Site ABCD",
    "domain":"siteabcd.com",
    "cat":["IAB2-1", "IAB2-2"],
    "privacypolicy":1,
    "page":"http://siteabcd.com/page.htm",
    "ref":"http://referringsite.com/referringpage.htm",
    "publisher":{
      "id":"pub12345",
      "name":"Publisher A"
    },
    "content":{
      "keywords":["keyword a","keyword b","keyword c"]
    }
  },
  "device":{
    "ip":"64.124.253.1",
    "ua":"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.6; en-US; rv:1.9.2.16) Gecko/20110319 Firefox/3.6.16",
    "os":"OS X",
    "flashver":"10.1",
    "js":1
  },
  "user":{
    "id":"45asdf987656789adfad4678rew656789",
    "buyeruid":"5df678asd8987656asdf78987654"
  }
}
`)

var expandableCreative []byte = []byte(`
{
  "id":"1234567893",
  "at":2,
  "tmax":120,
  "imp":[
    {
      "id":"1",
      "iframebuster":[
        "vendor1.com",
        "vendor2.com"
      ],
      "banner":{
        "w":300,
        "h":250,
        "pos":1,
        "battr":[
          13
        ],
        "expdir":[
          2,
          4
        ]
      }
    }
  ],
  "site":{
    "id":"1345135123",
    "name":"Site ABCD",
    "domain":"siteabcd.com",
    "sitecat":[
      "IAB2-1",
      "IAB2-2"
    ],
    "page":"http://siteabcd.com/page.htm",
    "ref":"http://referringsite.com/referringpage.htm",
    "privacypolicy":1,
    "publisher":{
      "id":"pub12345",
      "name":"Publisher A"
    },
    "content":{
      "keyword":[
        "keyword1",
        "keyword2",
        "keyword3"
      ]
    }
  },
  "device":{
    "ip":"64.124.253.1",
    "ua":"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.6; en-US; rv:1.9.2.16) Gecko/20110319 Firefox/3.6.16",
    "os":"OS X",
    "flashver":"10.1",
    "js":1
  },
  "user":{
    "id":"456789876567897654678987656789",
    "buyeruid":"545678765467876567898765678987654",
    "data":[
      {
        "id":"6",
        "name":"Data Provider 1",
        "segment":[
          {
            "id":"12341318394918",
            "name":"auto intenders"
          },
          {
            "id":"1234131839491234",
            "name":"auto enthusiasts"
          },
          {
            "id":"23423424",
            "name":"data-provider1-age",
            "value":"30-40"
          }
        ]
      }
    ]
  }
}`)
