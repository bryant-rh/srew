package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

//Write
func Write(data [][]string, out_type string, outFormat bool, mergeCell bool) {
	//var table *tablewriter.Table

	table := table(outFormat, mergeCell)
	var header []string
	//header = append(header, "NODE")
	switch {
	case out_type == "search":
		header = append(header,
			"NAME", "DESCRIPTION", "LATEST_VERSION", "INSTALLED", "UPGRADE",
		)
	case out_type == "search_all_version":
		header = append(header,
			"NAME", "DESCRIPTION", "VERSION", "INSTALLED",
		)
	case out_type == "list":
		header = append(header,
			"PLUGIN", "VERSION",
		)
	case out_type == "label":
		header = append(header,
			"Project Name", "Cluster Name", "Labels",
		)
	case out_type == "token":
		header = append(header,
			"SaToken",
		)
	default:
		header = append(header,
			"ID", "Project_ID", "Project Name", "Description", "CreatedAt", "UpdatedAt",
		)
	}
	table.SetHeader(header)
	// for _, i := range data {
	// 	table.Append(i)
	// }
	table.AppendBulk(data)

	table.Render()

}

//table
func table(outFormat, mergeCell bool) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	if outFormat {
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
		return table
	} else {
		if mergeCell {
			table.SetAutoMergeCells(true)

		}
		table.SetRowLine(true)
	}
	return table
}
