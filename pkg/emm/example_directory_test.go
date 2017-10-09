package emm

import "fmt"

func ExampleDirectory() {
	cd := `<directory>
		  <channel id="P_malekalssite">
			<dc:format>rss</dc:format>
			<dc:type>webnews</dc:type>
			<dc:subject>eucert</dc:subject>
			<dc:description>malekals site</dc:description>
			<dc:identifier>http://www.malekal.com/</dc:identifier>
			<iso:country>US</iso:country>
			<region>Global</region>
			<category>Specialist</category>
			<ranking>1</ranking>
			<iso:language>en</iso:language>
			<ocs:schedule>
			  <ocs:updatePeriod>daily</ocs:updatePeriod>
			  <ocs:updateFrequency>2</ocs:updateFrequency>
			</ocs:schedule>
			<feed title="malekals site" url="http://www.malekal.com/feed/"/>
		  </channel>
		</directory>`

	d := NewDirectory(cd)
	fmt.Println(d.Channels[0].ID)
	// Output:
	// P_malekalssite
}
