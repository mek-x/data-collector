package modules

import (
	/* collectors */
	_ "gitlab.com/mek_x/data-collector/pkg/collector/file"
	_ "gitlab.com/mek_x/data-collector/pkg/collector/mqtt"
	_ "gitlab.com/mek_x/data-collector/pkg/collector/shell"

	/* parsers */
	_ "gitlab.com/mek_x/data-collector/pkg/parser/jsonpath"
	_ "gitlab.com/mek_x/data-collector/pkg/parser/regex"

	/* dispatchers */
	_ "gitlab.com/mek_x/data-collector/pkg/dispatcher/cron"
	_ "gitlab.com/mek_x/data-collector/pkg/dispatcher/event"

	/* sinks */
	_ "gitlab.com/mek_x/data-collector/pkg/sink/gotify"
	_ "gitlab.com/mek_x/data-collector/pkg/sink/iotplotter"
	_ "gitlab.com/mek_x/data-collector/pkg/sink/stdout"
	_ "gitlab.com/mek_x/data-collector/pkg/sink/windy"
)
