package com.github.kr328.clash

import com.github.kr328.clash.common.util.intent
import com.github.kr328.clash.design.MainDesign
import com.github.kr328.clash.design.SettingsDesign
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.isActive
import kotlinx.coroutines.selects.select
import kotlinx.coroutines.withContext

class SettingsActivity : BaseActivity<SettingsDesign>() {
    override suspend fun main() {
        val design = SettingsDesign(this)

        setContentDesign(design)

        while (isActive) {
            select<Unit> {
                events.onReceive {

                }
                design.requests.onReceive {
                    when (it) {
                        SettingsDesign.Request.StartApp ->
                            startActivity(AppSettingsActivity::class.intent)
                        SettingsDesign.Request.StartNetwork ->
                            startActivity(NetworkSettingsActivity::class.intent)
                        SettingsDesign.Request.StartOverride ->
                            startActivity(OverrideSettingsActivity::class.intent)
                        SettingsDesign.Request.StartMetaFeature ->
                            startActivity(MetaFeatureSettingsActivity::class.intent)
                        SettingsDesign.Request.OpenHelp ->
                            startActivity(HelpActivity::class.intent)
                        SettingsDesign.Request.OpenAbout ->
                            design.showAbout(queryAppVersionName())
                        SettingsDesign.Request.OpenLogs ->
                            startActivity(LogsActivity::class.intent)
                    }
                }
            }
        }
    }
    private suspend fun queryAppVersionName(): String {
        return withContext(Dispatchers.IO) {
            packageManager.getPackageInfo(packageName, 0).versionName
        }
    }
}